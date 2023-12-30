package formutils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/go-playground/validator/v10"
	"github.com/monoculum/formam/v3"
)

// Convenient all-in-one function to parse, decode, and validate a POST request form.
//
// dst struct will be populated with the correct fields if all parsing, decoding and validating succeeds. Otherwise,
// this function returns an error and a HTTP response with the appropriate headers will be sent automatically.
//
// In case of validation errors, the error message for each field will be returned as a JSON object in the HTTP response.
func ParseDecodeValidate[T interface{}](app *app.App, w http.ResponseWriter, r *http.Request, dst *T) error {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Can't parse form", http.StatusBadRequest)
		return err
	}

	err = formam.Decode(r.PostForm, dst)
	if err != nil {
		obj, _ := err.(*formam.Error).MarshalJSON()
		http.Error(w, string(obj), http.StatusBadRequest)
		return err
	}

	err = app.Validate.Struct(dst)
	if err != nil {
		var ve validator.ValidationErrors
		veMap := make(map[string]string)
		if errors.As(err, &ve) {
			for _, fe := range ve {
				err := errorTagFunc[T](dst, fe.StructNamespace(), fe.Field(), fe.Tag())
				if err != nil {
					veMap[utils.ToLowerFirst(fe.Field())] = err.Error()
				}
			}
			jsonutils.Write(w, veMap, http.StatusBadRequest)
			return err
		}
		log.Panicln("Not a validator.ValidationErrors")
	}

	if r.URL.Query().Has("validate") {
		w.WriteHeader(http.StatusOK)
		return errors.New("Validate only")
	}

	return nil
}

func errorTagFunc[T interface{}](obj *T, snp string, fieldname, actualTag string) error {
	if !strings.Contains(snp, fieldname) {
		return nil
	}

	fieldArr := strings.Split(snp, ".")
	rsf := reflect.TypeOf(*obj)

	for i := 1; i < len(fieldArr); i++ {
		field, found := rsf.FieldByName(fieldArr[i])
		if found {
			if fieldArr[i] == fieldname {
				errorKey := field.Tag.Get("key")
				var message string
				if actualTag == "required" {
					message = "required.message"
				} else if actualTag == "min" {
					message = "tooShort.message"
				} else if actualTag == "max" {
					message = "tooLong.message"
				} else if actualTag == "username" {
					message = "invalidCharacters.message"
				} else if actualTag == "number" {
					message = "numbersOnly.message"
				} else {
					message = "invalid.message"
				}

				if errorKey != "" {
					return fmt.Errorf("%s.%s", errorKey, message)
				}
				return fmt.Errorf("%s", message)
			} else {
				if field.Type.Kind() == reflect.Ptr {
					// If the field type is a pointer, dereference it
					rsf = field.Type.Elem()
				} else {
					rsf = field.Type
				}
			}
		}
	}
	return nil
}
