package sourcefile

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/krostar/config/trivialerr"
)

func TestNewSource(t *testing.T) {
	var tests = map[string]struct {
		path              string
		opts              []Option
		expectedPath      string
		expectedExtension string
	}{
		"empty path": {
			path:              "",
			expectedPath:      "",
			expectedExtension: "",
		}, "yaml file": {
			path:              "path.yaml",
			expectedPath:      "path.yaml",
			expectedExtension: "yaml",
		}, "yml file": {
			path:              "path.yml",
			expectedPath:      "path.yml",
			expectedExtension: "yaml",
		}, "opts are applied": {
			path: "replace.me",
			opts: []Option{
				func(f *File) { f.path = "replaced" },
				func(f *File) { f.ext = "" },
			},
			expectedPath: "replaced",
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			file := NewSource(test.path, test.opts...)
			assert.Equal(t, test.expectedPath, file.path)
			assert.Equal(t, test.expectedExtension, file.ext)
			assert.True(t, file.strictOpen)
			assert.False(t, file.strictUnmarshal)
		})
	}
}

func TestFile_Unmarshal(t *testing.T) {
	type helloWorld struct {
		Hello string `json:"hello" yaml:"hello"`
	}
	var tests = map[string]struct {
		fileName               string
		fileContent            string
		createFile             bool
		expectedFailure        bool
		expectedTrivialFailure bool
		ffOpts                 []Option
		expectedTo             helloWorld
	}{
		"file cannot be open, without lazy opts": {
			fileName:               "file.toto",
			expectedFailure:        true,
			expectedTrivialFailure: false,
		}, "file cannot be open, with lazy opts": {
			fileName:               "file.toto",
			ffOpts:                 []Option{MayNotExist()},
			expectedFailure:        true,
			expectedTrivialFailure: true,
		}, "bli file": {
			createFile:             true,
			fileName:               "file.bli",
			expectedFailure:        true,
			expectedTrivialFailure: false,
		}, "json file": {
			createFile:  true,
			fileName:    "file.json",
			fileContent: `{"hello": "world"}`,
			expectedTo: helloWorld{
				Hello: "world",
			},
		}, "strict json file": {
			createFile:      true,
			fileName:        "file.json",
			fileContent:     `{"hello": "world", "world": "hello"}`,
			ffOpts:          []Option{FailOnUnknownFields()},
			expectedFailure: true,
		}, "yaml file": {
			createFile:  true,
			fileName:    "file.yaml",
			fileContent: `hello: "world"`,
			expectedTo: helloWorld{
				Hello: "world",
			},
		}, "strict yaml file": {
			createFile:      true,
			fileName:        "file.yaml",
			fileContent:     fmt.Sprintf("hello: \"world\"\nworld: \"hello\""),
			ffOpts:          []Option{FailOnUnknownFields()},
			expectedFailure: true,
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				fs   = afero.NewMemMapFs()
				opts = append(test.ffOpts, func(f *File) { f.fs = fs })
				file = NewSource(test.fileName, opts...)
				to   helloWorld
			)

			if test.createFile {
				require.NoError(t, afero.WriteFile(fs,
					test.fileName,
					[]byte(test.fileContent),
					0400,
				))
			}

			err := file.Unmarshal(&to)
			if test.expectedFailure {
				require.Error(t, err)
				assert.Equal(t, test.expectedTrivialFailure, trivialerr.IsTrivial(err))
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedTo, to)
			}
		})
	}
}

func TestFile_Name(t *testing.T) {
	require.Equal(t, "file", NewSource("").Name())
}
