package configue

// Source defines the interface any source should implements.
type Source interface {
	Name() string
}

type sourceThatUseReflection interface {
	Source
	GetReprValueByKey(key string) ([]byte, error)
}

type sourceThatUnmarshal interface {
	Source
	Unmarshal(cfg interface{}) error
}
