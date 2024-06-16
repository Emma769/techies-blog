package validator

import (
	"fmt"
	"strings"
)

type ValidationErrors map[string]string

func (e ValidationErrors) Error() string {
	var builder strings.Builder

	for k, v := range e {
		builder.Write([]byte(fmt.Sprintln(k, v)))
	}

	return builder.String()
}

type validator struct {
	*ValidationErrors
}

func New() *validator {
	return &validator{
		ValidationErrors: &ValidationErrors{},
	}
}

func (v validator) valid() bool {
	return len(*v.ValidationErrors) == 0
}

func (v *validator) add(key, msg string) {
	if _, ok := (*v.ValidationErrors)[key]; !ok {
		(*v.ValidationErrors)[key] = msg
	}
}

type validationfn[T any] func(T) (string, bool)

func Check[T any](v *validator, in T, fns ...validationfn[T]) error {
	for i := range len(fns) {
		errmsg, ok := fns[i](in)
		if !ok {
			parts := strings.Split(errmsg, ":")

			if len(parts) < 2 {
				panic("invalid error message - correct format \"key:message\"")
			}

			key, msg := parts[0], parts[1]

			v.add(key, msg)
		}
	}

	if !v.valid() {
		return v.ValidationErrors
	}

	return nil
}
