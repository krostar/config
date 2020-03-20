package config

import (
	"fmt"
	"strings"
)

// ValidateError contains all validation errors.
// Keys are field in error (with arboresence) and value the error.
// Key can be empty if the root interface failed.
type ValidateError map[string]error

// String implements Stringer for ValidateError.
func (v ValidateError) String() string {
	if len(v) == 0 {
		return "no validation errors"
	}

	var errors []string
	for key, value := range v {
		if key == "" {
			errors = append(errors, value.Error())
		} else {
			errors = append(errors, fmt.Sprintf("field %s %s", key, value.Error()))
		}
	}

	return "validation error: " + strings.Join(errors, ", ")
}

// Error implements Error for ValidateError.
func (v ValidateError) Error() string {
	return v.String()
}
