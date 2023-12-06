package formutils

import (
	"net/http"

	"github.com/doorman2137/betonz-go/internal/app"
	"github.com/go-playground/validator/v10"
	"github.com/monoculum/formam/v3"
)

func ParseDecodeValidate(app *app.App, w http.ResponseWriter, r *http.Request, dst interface{}) error {
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
		http.Error(w, err.(validator.ValidationErrors).Error(), http.StatusBadRequest)
		return err
	}

	return nil
}
