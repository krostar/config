package defaulter

import (
	"errors"
	"reflect"
)

type setDefaultFunc interface {
	SetDefault()
}

// SetDefault tries to set default to the provided interface.
// For a default to be applied, the interface needs to implement
// the `setDefaultFunc` interface, or to contain any fields (in
// case of a struct) that implement it. The SetDefault method will
// be called only if the value is the zero value.
func SetDefault(i interface{}) error {
	var value = reflect.ValueOf(i)

	// if the reflected value is nil
	if !value.IsValid() {
		return errors.New("i value is nil")
	}

	// as i is actually an interface, just get the thing behind the interface
	value = reflect.Indirect(value.Elem())

	// recursively walk through it to try to set default
	walkThrough(&value)

	return nil
}

func walkThrough(v *reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr:
		// we don't want to put default on nil pointed value so just leave here
		if v.IsNil() {
			break
		}

		// retrieve the pointed value
		pv := v.Elem()

		// try to put a default recursively to the pointed value
		walkThrough(&pv)
	case reflect.Struct:
		const tagKey = "default"
		// try to put default on the whole structure
		tryToSetDefault(v)

		// try each fields, there may be empty fields that also want defaults
		for i := 0; i < v.NumField(); i++ {
			// get field
			childV := v.Field(i)
			childField := v.Type().Field(i)

			// check if the tag discard the field
			if tag := childField.Tag.Get(tagKey); tag == "-" {
				continue
			}

			// ignore it if the field is unexported
			if v.Type().Field(i).PkgPath != "" {
				continue
			}

			// handle the child recursively
			walkThrough(&childV)
		}
	default:
		// for every other types, try to set default
		tryToSetDefault(v)
	}
}

func tryToSetDefault(v *reflect.Value) {
	// since we're gonna need the value to match an interface through the
	// address of the struct, make sure we can do that
	if !v.IsValid() || !v.CanInterface() || !v.CanAddr() {
		return
	}

	// since the defaulter updates itself, it needs to receive a pointor;
	// so try if the value's pointor's interface implements the defaulter interface
	if f, ok := v.Addr().Interface().(setDefaultFunc); ok {
		// if the value is the zero value, call the defaulter
		if isZeroValue(v) {
			f.SetDefault()
		}
	}
}

func isZeroValue(v *reflect.Value) bool {
	if !v.IsValid() || !v.CanInterface() {
		return false
	}
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
