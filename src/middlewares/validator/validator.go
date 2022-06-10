package validator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	once     sync.Once
	instance *validator.Validate
)

func Validate() *validator.Validate {
	once.Do(func() {
		instance = validator.New()

		// register function to get tag name from json tags
		instance.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	})

	return instance
}

func ReportError(err error) error {
	message := "`%s` is invalid: %s %s"
	message = fmt.Sprintf(message,
		err.(validator.ValidationErrors)[0].Field(),
		err.(validator.ValidationErrors)[0].ActualTag(),
		err.(validator.ValidationErrors)[0].Param())

	return fmt.Errorf(strings.TrimSpace(message))
}
