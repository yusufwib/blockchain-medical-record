// mvalidator wrapper dari validator, menambahkan translator dan error response dalam bentuk map[string]
package mvalidator

import (
	"errors"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator interface {
	Struct(input interface{}) (map[string]interface{}, error)
}

type mValidator struct {
	instance   *validator.Validate
	translator ut.Translator
}

func New() *mValidator {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Panic("translator not found")
	}

	_ = enTranslations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if len(name) == 0 {
			name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})

	_ = validate.RegisterTranslation("ddate", trans, func(ut ut.Translator) error {
		return ut.Add("ddate", "{0} must be valid date format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ddate", fe.Field())
		return t
	})

	_ = validate.RegisterValidation("ddate", func(fl validator.FieldLevel) bool {
		str := fl.Field().String()
		layout := "2006-01-02 15:04:05"
		_, err := time.Parse(layout, str)
		return err == nil
	})

	_ = validate.RegisterTranslation("cdate", trans, func(ut ut.Translator) error {
		return ut.Add("cdate", "{0} must be valid date format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("cdate", fe.Field())
		return t
	})

	_ = validate.RegisterValidation("cdate", func(fl validator.FieldLevel) bool {
		str := fl.Field().String()
		if str == "" {
			return true
		}
		layout := "2006-01-02"
		_, err := time.Parse(layout, str)
		return err == nil
	})

	return &mValidator{
		instance:   validate,
		translator: trans,
	}
}

func (m *mValidator) Struct(input interface{}) (map[string]interface{}, error) {

	errMap := make(map[string]interface{})

	if input == nil {
		errorMsg := "developer_vault. input cant be nil"
		errMap["message"] = errorMsg
		return errMap, errors.New(errorMsg)
	}

	err := m.instance.Struct(input)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			errMap["message"] = err.Error()
			return errMap, err
		}

		errs := err.(validator.ValidationErrors)
		for _, e := range errs {

			structName := strings.Split(e.Namespace(), ".")[0]
			jsonFieldName := strings.Replace(e.Namespace(), structName+".", "", 1)
			errMap[jsonFieldName] = e.Translate(m.translator)
		}
		return errMap, errs
	}

	return nil, nil
}
