package validator

import (
    "fmt"
    "github.com/go-playground/locales/en"
    ut "github.com/go-playground/universal-translator"
    "github.com/go-playground/validator/v10"
    entranslations "github.com/go-playground/validator/v10/translations/en"
    "log"
    "net/http"
    "strings"
)

type humanReadableErrors struct {
    errList []string
}

func (err humanReadableErrors) Error() string {
    return strings.Join(err.errList, "\n")
}

func (err *humanReadableErrors) addError(errStr string) {
    err.errList = append(err.errList, errStr)
}

type CustomValidator struct {
    validate *validator.Validate
    trans    ut.Translator
}

func NewValidator() *CustomValidator {
    v := validator.New()

    translator := en.New()
    uni := ut.New(translator, translator)

    trans, _ := uni.GetTranslator("en")

    if err := entranslations.RegisterDefaultTranslations(v, trans); err != nil {
        log.Fatal(err)
    }

    // human readable error messages
    translations := []struct {
        tag         string
        translation string
    }{
        {
            tag:         "uri",
            translation: fmt.Sprintf("'{0}' is not valid URI"),
        },
        {
            tag:         "required",
            translation: fmt.Sprintf("'{0}' is required"),
        },
        {
            tag:         "httpStatus",
            translation: fmt.Sprintf("'{0}' is not a valid HTTP status code"),
        },
    }

    for _, t := range translations {
        err := v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation), translateFunc)
        if err != nil {
            panic(err)
        }
    }

    _ = v.RegisterValidation("httpStatus", httpStatusValidator)
    return &CustomValidator{
        validate: v,
        trans: trans,
    }
}

func httpStatusValidator(fl validator.FieldLevel) bool {
    return http.StatusText(int(fl.Field().Int())) != ""
}

func registrationFunc(tag string, translation string) validator.RegisterTranslationsFunc {
    return func(ut ut.Translator) (err error) {
        if err = ut.Add(tag, translation, true); err != nil {
            return
        }
        return
    }
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
    valueStr := fmt.Sprintf("%v", fe.Value())
    t, err := ut.T(fe.Tag(), valueStr)
    if err != nil {
        return fe.(error).Error()
    }
    return t
}

func (v *CustomValidator) ValidateStruct(value interface{}) error {
    // validate the given struct
    err := v.validate.Struct(value)
    if err == nil {
        return nil
    }

    result := humanReadableErrors{}
    vErrors := err.(validator.ValidationErrors)
    for _, vErr := range vErrors {
        result.addError(vErr.Translate(v.trans))
    }

    return result
}
