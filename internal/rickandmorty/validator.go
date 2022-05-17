package rickandmorty

import (
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var validate *validator.Validate
var translator ut.Translator

type validationErrors struct {
	Errors []string `json:"errors"`
}

func getField(field string) string {
	fieldParts := strings.Split(field, ".")
	if len(fieldParts) != 2 {
		return field
	}
	return fieldParts[1]
}

func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		return err
	}
	return nil
}

func Translate(err error) *validationErrors {
	if errs, ok := err.(validator.ValidationErrors); ok {
		vErr := []string{}
		errMap := errs.Translate(translator)
		for _, errString := range errMap {
			vErr = append(vErr, errString)
		}
		return &validationErrors{Errors: vErr}
	}
	return nil
}

func init() {
	en := en.New()
	uni := ut.New(en, en)
	translator, _ = uni.GetTranslator("en")
	validate = validator.New()
	en_translations.RegisterDefaultTranslations(validate, translator)
}
