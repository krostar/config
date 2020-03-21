package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/krostar/config/internal/trivialerr"
)

// InitializeNewValueOfTypeWithJSON .
// nolint: gocyclo
func InitializeNewValueOfTypeWithJSON(typ reflect.Type, jsonData []byte) (*reflect.Value, error) {
	if typ == nil {
		return nil, errors.New("cannot create value of nil type")
	}

	vPtr := reflect.New(typ)
	v := vPtr.Elem()

	var err error

	// use the json unmarshaller as it already handle the super fat switch
	// on types/values except for the case of a time.Duration, unmarshall
	// in the custom duration type and manually set the value; because json
	// unmarshaler does not handle time.Duration in string format.
	switch typ {
	case reflect.TypeOf(time.Second):
		var cd customDuration
		if err = json.Unmarshal(jsonData, &cd); err != nil {
			err = fmt.Errorf("custom duration unmarshal failed: %w", err)
			break
		}
		v.SetInt(cd.ToInt64())
	default:
		if err = json.Unmarshal(jsonData, vPtr.Interface()); err != nil {
			err = fmt.Errorf("json marshaller failed to fill the value: %w", err)
			break
		}
	}

	if err != nil {
		return nil, err
	}

	return &v, nil
}

// InitializeNewValueOfTypeWithString .
func InitializeNewValueOfTypeWithString(typ reflect.Type, str string) (*reflect.Value, error) {
	if typ == nil {
		return nil, errors.New("cannot create value of nil type")
	}

	kind := typ.Kind()
	underTyp := typ

	switch kind {
	// we may need to transform the representation to fit with json unmarshaller
	// so if it's a ptr, give us the pointed type
	case reflect.Ptr:
		underTyp = typ.Elem()
		kind = underTyp.Kind()
	// if we got a interface we can't guess the final type, put the value as a string
	case reflect.Interface:
		kind = reflect.String
	}

	// once we have a valid type, if it's a time.Duration, wrap it as json unmarshaller do not handle it
	if underTyp.Comparable() && underTyp == reflect.TypeOf(time.Second) {
		// time.Duration can be of written in float64 representation or in string
		kind = reflect.String
		if _, err := strconv.ParseFloat(str, 64); err == nil {
			kind = reflect.Float64
		}
	}

	// if the representation is a string, quote it as json unmarshaller need them
	if kind == reflect.String {
		str = strconv.Quote(str)
	}

	return InitializeNewValueOfTypeWithJSON(typ, []byte(str))
}

// SetNewValue .
func SetNewValue(oldValue, newValue *reflect.Value) (bool, error) {
	if !oldValue.CanSet() {
		return false, errors.New("value is not settable")
	}

	oldValue.Set(*newValue)

	return true, nil
}

func setValuesForEachAttributes(src SourceSetValueFromConfigTreePath, cfg interface{}) error {
	var value = reflect.ValueOf(cfg)

	if value.IsNil() {
		return errors.New("cfg is nil")
	}

	value = reflect.Indirect(value.Elem())
	if _, err := setValueRecursively(src, "", &value); err != nil {
		return err
	}

	return nil
}

func setValueRecursively(src SourceSetValueFromConfigTreePath, path string, v *reflect.Value) (bool, error) {
	switch v.Kind() {
	case reflect.Invalid:
		return false, errors.New("value is invalid")
	case reflect.Ptr:
		return setValuePointor(src, path, v)
	case reflect.Struct:
		return setValueStruct(src, path, v)
	default:
		isset, err := src.SetValueFromConfigTreePath(v, path)
		if err != nil {
			if !trivialerr.IsTrivial(err) {
				return false, fmt.Errorf("unable to get value for key %q: %w", path, err)
			}
			return false, nil
		}
		return isset, nil
	}
}

func setValuePointor(src SourceSetValueFromConfigTreePath, path string, v *reflect.Value) (bool, error) {
	var validV = *v

	// if we have a nil pointor, build a non-nil one
	if v.IsNil() {
		validV = reflect.New(v.Type().Elem())
	}

	// retrieve the pointed value
	newV := validV.Elem()

	// go recursively with the pointed value
	if isSet, err := setValueRecursively(src, path, &newV); err != nil || !isSet {
		return false, err
	}

	// if a value has been set, use it (otherwise don't replace the original value)
	if !v.CanSet() {
		return false, errors.New("value is not settable")
	}

	v.Set(validV)
	return true, nil
}

func setValueStruct(src SourceSetValueFromConfigTreePath, path string, v *reflect.Value) (bool, error) {
	const tagKey = "cfg"

	var oneIsSet = false

	for i := 0; i < v.NumField(); i++ {
		// get info of each fields as for example the real name of the field
		var (
			childV     = v.Field(i)
			childField = v.Type().Field(i)
			childPath  = appendConfigTreePath(path, childField.Name)
		)

		// if a tag is defined, override the name with it
		if tag := childField.Tag.Get(tagKey); tag != "" {
			childPath = appendConfigTreePath(path, tag)
		}

		// ignore if the field is unexported or the name is `-`
		if childField.PkgPath != "" || path == "-" {
			continue
		}

		// recursive call with the value
		if isSet, err := setValueRecursively(src, childPath, &childV); err == nil {
			if isSet {
				oneIsSet = true
			}
		} else {
			return isSet, err
		}
	}

	return oneIsSet, nil
}
