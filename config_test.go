package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_success(t *testing.T) {
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

func TestLoad_failures(t *testing.T) {
	var cfg interface{}

	// nil receiver
	require.Error(t, Load(nil))
	// non-ptr receiver
	require.Error(t, Load(cfg))
	// unknown type of source
	require.Error(t, Load(&cfg, WithSources(&dumbSource{})))
	// source that returns no errors
	require.NoError(t, Load(&cfg, WithSources(stubSourceThatUnmarshal(0))))
	// source that returns a trivial error
	require.NoError(t, Load(&cfg, WithSources(stubSourceThatUnmarshal(-1))))
	// source that returns a real error
	require.Error(t, Load(&cfg, WithSources(stubSourceThatUnmarshal(1))))
}

func TestConfig_Load_opts_applied(t *testing.T) {
	var (
		cfg int
		s1  = dumbSource{}
		s2  = dumbSource{}
	)
	var c = New(WithSources(s2))

	// this is gonna fail as dumbSource does not really implement a usesable source
	assert.Error(t, c.Load(&cfg, WithSourcesPrepend(s1)))

	assert.Equal(t, []Source{s1, s2}, c.sources)
}
