package config

// Source defines the interface any source should implements.
type Source interface {
	Name() string
}

// SourceGetReprValueByKey defines the way a source should
// return a representative value based on a key.
type SourceGetReprValueByKey interface {
	Source
	GetReprValueByKey(key string) ([]byte, error)
}

// SourceUnmarshal defines a way to apply a source to
// a config via unmarshalling directly to it.
type SourceUnmarshal interface {
	Source
	Unmarshal(cfg interface{}) error
}
