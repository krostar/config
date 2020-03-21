package config

import (
	"errors"
	"reflect"
)

type validateFunc interface {
	Validate() error
}

// Validate calls validateFunc for all types that implements it.
func Validate(i interface{}) error {
	var value = reflect.ValueOf(i)

	// if the reflected value is nil
	if !value.IsValid() {
		return errors.New("i value is nil")
	}

	// as i is actually an interface, just get the thing behind the interface
	value = reflect.Indirect(value.Elem())

	var errs = make(ValidationError)
	// recursively walk through it to see if it's valid
	validateRecursively(&value, "", errs)

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateRecursively(v *reflect.Value, name string, errs ValidationError) {
	switch v.Kind() {
	case reflect.Ptr:
		// we don't want to put default on nil pointed value so just leave here
		if v.IsNil() {
			break
		}

		// retrieve the pointed value
		var pv = v.Elem()

		// try to put a default recursively to the pointed value
		validateRecursively(&pv, name, errs)
	case reflect.Struct:
		// try to put default on the whole structure
		if err := validateValue(v); err != nil {
			errs[name] = err
		}

		// try each fields, there may be empty fields that also want defaults
		for i := 0; i < v.NumField(); i++ {
			// get field
			var (
				childV     = v.Field(i)
				childField = v.Type().Field(i)
			)

			// ignore it if the field is unexported
			if v.Type().Field(i).PkgPath != "" {
				continue
			}

			// handle the child recursively
			validateRecursively(&childV, appendConfigTreePath(name, childField.Name), errs)
		}
	default:
		// for every other types, try to set default
		if err := validateValue(v); err != nil {
			errs[name] = err
		}
	}
}

func validateValue(v *reflect.Value) error {
	// since we're gonna need the value to match an interface through the
	// address of the struct, make sure we can do that
	if !v.IsValid() || !v.CanInterface() || !v.CanAddr() {
		return nil
	}

	if f, ok := v.Addr().Interface().(validateFunc); ok {
		return f.Validate()
	}

	return nil
}
