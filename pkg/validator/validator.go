package validator

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/gookit/slog"
)

type Validate struct {
	*validator.Validate
}

func NewValidator() *Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.RegisterValidation("rfc3339", ValidateRFC3339); err != nil {
		slog.Fatal("Failed to register rfc3339 validation", "error", err)
		os.Exit(1)
	}

	if err := validate.RegisterValidation("action_type", ValidateActionType); err != nil {
		slog.Fatal("Failed to register action_type validation", "error", err)
		os.Exit(1)
	}

	if err := validate.RegisterValidation("role", ValidateRole); err != nil {
		slog.Fatal("Failed to register role validation", "error", err)
		os.Exit(1)
	}

	if err := validate.RegisterValidation("letters_only", ValidateLettersOnly); err != nil {
		slog.Fatal("Failed to register letters_only validation", "error", err)
		os.Exit(1)
	}

	return &Validate{
		Validate: validate,
	}
}

func (v *Validate) FormatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			fieldName := e.Field()
			tag := e.Tag()
			msg := fmt.Sprintf("'%s' failed on the '%s' tag", fieldName, tag)
			return msg
		}
	}
	errorMsg := err.Error()
	errorMsg = strings.TrimPrefix(errorMsg, "validation for ")

	return errorMsg
}

func ValidateRFC3339(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, err := time.Parse(time.RFC3339, value)
	return err == nil
}

func ValidateActionType(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := ActionType(field.String())
	if _, ok := AllowedActionTypes[value]; ok {
		return true
	}

	return false
}

func ValidateRole(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	value := Role(field.String())
	if _, ok := AllowedRoles[value]; ok {
		return true
	}

	return false
}

func ValidateLettersOnly(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	for _, r := range value {
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}
