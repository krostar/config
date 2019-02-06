package configue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	type (
		icfgNested struct {
			withUnexportedField int
			WithCustomTag       string `cfg:"olleh"`
			UniversalStrAnswer  string
		}

		icfg struct {
			WithDiscardTag            string `cfg:"-"`
			WithDiscardTagWithDefault string `cfg:"-"`
			WithCustomTag             string `cfg:"hello"`
			UniversalAnswer           int
			withUnexportedField       int
			WithDefaultValue          string
			WithUntouchedDefaultValue string
			WithNestedStruct          icfgNested
			WithInterface             interface{}
			WithNilPointer            *string
			WithPointer               *int
		}
	)
	var (
		i      = 42
		str    = "I'm not nil anymore"
		source = stubSourceThatUseReflection{
			"hello":                                []byte("world"),
			"universalanswer":                      []byte("42"),
			"withunexportedfield":                  []byte("-30"),
			"withdefaultvalue":                     []byte("I've been replaced"),
			"withnestedstruct.withunexportedfield": []byte("-30"),
			"withnestedstruct.olleh":               []byte("dlrow"),
			"withnestedstruct.universalstranswer":  []byte("forty-two"),
			"withinterface":                        []byte("blih"),
			"withnilpointer":                       []byte("I'm not nil anymore"),
			"withpointer":                          []byte("42"),
		}
		cfg = icfg{
			WithDiscardTagWithDefault: "I should NOT be replaced",
			withUnexportedField:       30,
			WithDefaultValue:          "I should be replaced",
			WithUntouchedDefaultValue: "I should NOT be replaced",
			WithPointer:               new(int),
		}
		expectedCfg = icfg{
			WithDiscardTag:            "",
			WithDiscardTagWithDefault: "I should NOT be replaced",
			WithCustomTag:             "world",
			UniversalAnswer:           42,
			withUnexportedField:       30,
			WithDefaultValue:          "I've been replaced",
			WithUntouchedDefaultValue: "I should NOT be replaced",
			WithNestedStruct: icfgNested{
				withUnexportedField: 0,
				WithCustomTag:       "dlrow",
				UniversalStrAnswer:  "forty-two",
			},
			WithInterface:  "blih",
			WithNilPointer: &str,
			WithPointer:    &i,
		}
	)
	err := Load(&cfg, WithSources(source))
	require.NoError(t, err)
	assert.Equal(t, expectedCfg, cfg)
}
