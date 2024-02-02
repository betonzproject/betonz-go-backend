package formutils

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/doorman2137/betonz-go/internal/utils"
	"github.com/doorman2137/betonz-go/internal/utils/fileutils"
	"github.com/doorman2137/betonz-go/internal/utils/jsonutils"
	"github.com/go-playground/validator/v10"
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

	_ = app.Decoder.Decode(dst, r.PostForm)
	err = app.Validate.Struct(dst)
	if err != nil {
		var ve validator.ValidationErrors
		veMap := make(map[string]string)
		if errors.As(err, &ve) {
			for _, fe := range ve {
				err := errorTagFunc[T](dst, fe.Field(), fe.Tag())
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

func ParseDecodeValidateMultipart[T interface{}](app *app.App, w http.ResponseWriter, r *http.Request, dst *T) error {
	multipart, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Can't parse multipart", http.StatusBadRequest)
		return err
	}

	values := make(url.Values)
	tempFileSums := make(map[string]string)
	errorMap := make(map[string]string)
	hasError := false
	for {
		part, err := multipart.NextPart()
		if err == io.EOF {
			break
		}
		defer part.Close()

		if part.FileName() == "" {
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(part)
			if err != nil {
				http.Error(w, "Can't parse multipart", http.StatusBadRequest)
			}
			values.Add(part.FormName(), buf.String())
		} else if !r.URL.Query().Has("validate") {
			tempFile, err := os.CreateTemp("", part.FormName())
			if err != nil {
				log.Panicln("Can't create temp directory: " + err.Error())
			}
			defer fileutils.CloseAndDelete(tempFile)

			hash := sha256.New()

			multiWriter := io.MultiWriter(tempFile, hash)

			_, err = io.Copy(multiWriter, part)
			if err != nil {
				log.Panicln(err)
			}

			sum := fmt.Sprintf("%x", hash.Sum(nil))
			tempFileSums[tempFile.Name()] = sum

			if !fileutils.IsSupportedFileType(tempFile.Name()) {
				errorMap[part.FormName()] = "upload.fileTypeNotSupported.message"
				hasError = true
			}

			values.Add(part.FormName(), sum)
		}
	}

	_ = app.Decoder.Decode(dst, values)
	err = app.Validate.Struct(dst)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				err := errorTagFunc[T](dst, fe.Field(), fe.Tag())
				if err != nil && errorMap[fe.Field()] == "" {
					errorMap[utils.ToLowerFirst(fe.Field())] = err.Error()
					hasError = true
				}
			}
		} else {
			log.Panicln("Not a validator.ValidationErrors")
		}
	}

	if hasError {
		jsonutils.Write(w, errorMap, http.StatusBadRequest)
		return errors.New("Upload error")
	}

	if r.URL.Query().Has("validate") {
		w.WriteHeader(http.StatusOK)
		return errors.New("Validate only")
	}

	// Copy uploaded files to uploads folder
	for filename, sum := range tempFileSums {
		path := filepath.Join("uploads", sum[:2], sum[2:4], sum[4:])
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		err = fileutils.Copy(filename, path)
		if err != nil {
			return err
		}
	}

	return nil
}

func errorTagFunc[T interface{}](obj *T, fieldname, actualTag string) error {
	t := reflect.TypeOf(*obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, found := t.FieldByName(fieldname)
	if found {
		errorKey := field.Tag.Get("key")
		var message string

		switch actualTag {
		case "required":
			message = "required.message"
		case "min":
			if field.Type.Kind() == reflect.Int64 {
				message = "tooLow.message"
			} else {
				message = "tooShort.message"
			}
		case "max":
			if field.Type.Kind() == reflect.Int64 {
				message = "tooHigh.message"
			} else {
				message = "tooLong.message"
			}
		case "username":
			message = "invalidCharacters.message"
		case "number":
			message = "numbersOnly.message"
		default:
			message = "invalid.message"
		}

		if errorKey != "" {
			return fmt.Errorf("%s.%s", errorKey, message)
		}
		return fmt.Errorf("%s", message)
	}

	return nil
}
