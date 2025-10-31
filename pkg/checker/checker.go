package checker

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslators "github.com/go-playground/validator/v10/translations/en"
	"github.com/quantsmithapp/datastation-backend/pkg/errors"
)

type checker struct {
	validate   *validator.Validate
	translator ut.Translator
}

func NewChecker() *checker {
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, found := uni.GetTranslator("en")
	if !found {
		panic("locale language not found")
	}

	_ = enTranslators.RegisterDefaultTranslations(validate, trans)
	return &checker{validate: validate, translator: trans}
}

func (c checker) Struct(s interface{}) error {
	err := c.validate.Struct(s)
	return c.translate(err)
}

func (c checker) Var(field interface{}, tag string) error {
	err := c.validate.Var(field, tag)
	return c.translate(err)
}

func (c checker) translate(err error) error {
	if err == nil {
		return nil
	}

	return c.translates(err)[0]
}

func (c checker) translates(err error) []error {
	if err == nil {
		return nil
	}

	validatorErrs := err.(validator.ValidationErrors)
	var errs []error
	for _, e := range validatorErrs {
		parsed := errors.NewBadRequest("VAL400", e.Translate(c.translator))
		errs = append(errs, parsed)
	}

	return errs
}
