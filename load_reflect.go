package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/krostar/config/trivialerr"
)

func loadThroughReflection(source SourceGetReprValueByKey, cfg interface{}) error {
	var value = reflect.ValueOf(cfg)

	if value.IsNil() {
		return errors.New("cfg is nil")
	}

	value = reflect.Indirect(value.Elem())
	if _, err := loadReflectRecursivly(source, "", &value); err != nil {
		return err
	}

	return nil
}

func loadReflectRecursivly(source SourceGetReprValueByKey, name string, v *reflect.Value) (bool, error) {
	switch v.Kind() {
	case reflect.Invalid:
		return false, errors.New("value is invalid")
	case reflect.Ptr:
		return loadReflectHandlePointer(source, name, v)
	case reflect.Struct:
		return loadReflectHandleStruct(source, name, v)
	default:
		return loadReflectHandleDefault(source, name, v)
	}
}

func loadReflectHandlePointer(source SourceGetReprValueByKey, name string, v *reflect.Value) (bool, error) {
	var validV = *v

	// if we have a nil pointor, build a non-nil one
	if v.IsNil() {
		validV = reflect.New(v.Type().Elem())
	}

	// retrieve the pointed value
	newV := validV.Elem()

	// go recursively with the pointed value
	if isSet, err := loadReflectRecursivly(source, name, &newV); err != nil || !isSet {
		return false, err
	}

	// if a value has been set, use it (otherwise don't replace the original value)
	if !v.CanSet() {
		return false, errors.New("value is not settable")
	}

	v.Set(validV)
	return true, nil
}

func loadReflectHandleStruct(source SourceGetReprValueByKey, name string, v *reflect.Value) (bool, error) {
	const tagKey = "cfg"
	var oneIsSet = false

	for i := 0; i < v.NumField(); i++ {
		// get info of each fields as for example the real name of the field
		var (
			childV     = v.Field(i)
			childField = v.Type().Field(i)
			childName  = fieldNamer(name, childField.Name)
		)

		// if a tag is defined, override the name with it
		if tag := childField.Tag.Get(tagKey); tag != "" {
			childName = fieldNamer(name, tag)
		}

		// ignore if the field is unexported or the name is `-`
		if childField.PkgPath != "" || name == "-" {
			continue
		}

		// recursive call with the value
		if isSet, err := loadReflectRecursivly(source, childName, &childV); err == nil {
			if isSet {
				oneIsSet = true
			}
		} else {
			return isSet, err
		}
	}

	return oneIsSet, nil
}

func loadReflectHandleDefault(source SourceGetReprValueByKey, name string, v *reflect.Value) (bool, error) {
	// asks nicely if source has a value for this key
	repr, err := source.GetReprValueByKey(name)

	// if source failed, leave
	if err != nil {
		if !trivialerr.IsTrivial(err) {
			return false, fmt.Errorf("unable to get value for key %q: %w", name, err)
		}
		return false, nil
	}

	// otherwise try to use the value found
	var newValue *reflect.Value
	if newValue, err = createNewValueOfType(repr, v.Type()); err != nil {
		return false, fmt.Errorf("unable to convert value for key %q: %w", name, err)
	}

	if !v.CanSet() {
		return false, errors.New("value is not settable")
	}

	v.Set(*newValue)
	return true, nil
}

func fieldNamer(parentName string, childName string) string {
	if parentName != "" {
		childName = parentName + "." + childName
	}
	return strings.ToLower(childName)
}

// nolint: gocyclo
func createNewValueOfType(repr []byte, typ reflect.Type) (*reflect.Value, error) {
	if typ == nil {
		return nil, errors.New("cannot create value of nil type")
	}

	var (
		// create a new value with the wanted type
		vPtr = reflect.New(typ)
		v    = vPtr.Elem()
		kind = typ.Kind()
	)

	// we may need to transform the representation to fit with json unmarshaller
	// so if it's a ptr, give us the pointed type
	if kind == reflect.Ptr {
		typ = typ.Elem()
		kind = typ.Kind()
	}

	// if we got a interface we can't guess the final type, put the value as a string
	if kind == reflect.Interface {
		kind = reflect.String
	}

	// once we have a valid type, if it's a time.Duration, wrap it as json
	// unmarshaller do not handle it
	if typ.Comparable() && typ == reflect.TypeOf(time.Second) {
		// time.Duration can be of written in float64 representation or in string
		kind = reflect.String
		if _, err := strconv.ParseFloat(string(repr), 64); err == nil {
			kind = reflect.Float64
		}
	}

	// if the representation is a string, quote it as json unmarshaller need them
	if kind == reflect.String {
		repr = []byte(strconv.Quote(string(repr)))
	}

	// use the json unmarshaller as it already handle the super fat switch on types/values
	var sErr error
	switch typ {
	case reflect.TypeOf(time.Second):
		// in case of a time.Duration, unmarshall in the custom duration type
		// and manually set the value.
		var cd customDuration
		if err := json.Unmarshal(repr, &cd); err != nil {
			sErr = fmt.Errorf("custom duration unmarshal failed: %w", err)
			break
		}
		v.SetInt(cd.ToInt64())
	default:
		if err := json.Unmarshal(repr, vPtr.Interface()); err != nil {
			sErr = fmt.Errorf("json marshaller failed to fill the value: %w", err)
			break
		}
	}

	if sErr != nil {
		return nil, sErr
	}

	return &v, nil
}
