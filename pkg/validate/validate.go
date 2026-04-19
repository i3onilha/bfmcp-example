// Package validate configures go-playground/validator with project-wide custom tags.
package validate

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// New returns a validator instance with custom tags such as "notblank".
func New() (*validator.Validate, error) {
	v := validator.New()
	if err := v.RegisterValidation("notblank", notBlank); err != nil {
		return nil, fmt.Errorf("register notblank: %w", err)
	}
	return v, nil
}

func notBlank(fl validator.FieldLevel) bool {
	s, ok := fl.Field().Interface().(string)
	if !ok {
		return true
	}
	return strings.TrimSpace(s) != ""
}
