package config

import (
	"reflect"
	"strings"
)

// SourceCreationFunc defines how to  creates a new source.
type SourceCreationFunc func() (Source, error)

// Source defines the interface any source must implements.
type Source interface {
	Name() string
}

// SourceUnmarshal defines a way to apply a source to
// a config via unmarshalling directly to it.
type SourceUnmarshal interface {
	Source
	Unmarshal(cfg interface{}) error
}

// SourceSetValueFromConfigTreePath defines a way to set explicitly
// each values from each configuration paths.
type SourceSetValueFromConfigTreePath interface {
	Source
	SetValueFromConfigTreePath(v *reflect.Value, treePath string) (bool, error)
}

func appendConfigTreePath(parentPath string, childName string) string {
	if parentPath != "" {
		childName = parentPath + "." + childName
	}
	return strings.ToLower(childName)
}
