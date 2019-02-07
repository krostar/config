package sourcefile

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"

	"github.com/krostar/config/trivialerr"
)

// File implements config.Source to fetch values from a file.
type File struct {
	fs              afero.Fs
	path            string
	ext             string
	strictUnmarshal bool
	strictOpen      bool
}

// NewSource returns a new file source.
func NewSource(path string, opts ...Option) *File {
	var ext = filepath.Ext(path)
	if ext != "" {
		ext = strings.ToLower(ext[1:])
	}

	if ext == "yml" {
		ext = "yaml"
	}

	ff := File{
		fs:              afero.NewReadOnlyFs(afero.NewOsFs()),
		path:            path,
		ext:             ext,
		strictUnmarshal: false,
		strictOpen:      true,
	}

	for _, opt := range opts {
		opt(&ff)
	}

	return &ff
}

// Name implements config.Source interface.
func (f *File) Name() string {
	return "file"
}

// Unmarshal tries to unmarshal file to the provided interface.
// It returns a trivial error if load strictness is false, or the true error otherwise.
func (f *File) Unmarshal(to interface{}) error {
	ff, err := f.fs.Open(f.path)
	if err != nil {
		return trivialerr.WrapIf(f.strictOpen,
			errors.Wrapf(err, "unable to open file %q", f.path),
		)
	}
	defer ff.Close() // nolint: errcheck, gosec

	switch f.ext {
	case "json":
		var decoder = json.NewDecoder(ff)
		if f.strictUnmarshal {
			decoder.DisallowUnknownFields()
		}
		err = decoder.Decode(to)
	case "yaml":
		var decoder = yaml.NewDecoder(ff)
		if f.strictUnmarshal {
			decoder.SetStrict(f.strictUnmarshal)
		}
		err = decoder.Decode(to)
	default:
		err = errors.Errorf("%q extension is not supported", f.ext)
	}

	return errors.Wrapf(err, "failed to unmarshal file %q", f.path)
}
