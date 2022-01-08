package validate

import (
	"reflect"
	"strings"

	pt_br "github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	pt_br_translations "github.com/go-playground/validator/v10/translations/pt_BR"
	"github.com/google/uuid"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator ut.Translator

func init() {
	validate = validator.New()

	translator, _ = ut.New(pt_br.New(), pt_br.New()).GetTranslator("pt_br")

	pt_br_translations.RegisterDefaultTranslations(validate, translator)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Check(val interface{}) error {
	if err := validate.Struct(val); err != nil {

		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		var fields FieldErrors
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(translator),
			}
			fields = append(fields, field)
		}

		return fields
	}

	return nil
}

func GenerateID() string {
	return uuid.NewString()
}

func CheckID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}
	return nil
}
