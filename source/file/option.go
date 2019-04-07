package sourcefile

// Option defines the function signature to apply options.
type Option func(f *File)

// MayNotExist tells the file.Unmarshal function to return a
// that implements IsTrivial when file does not exists.
func MayNotExist() Option {
	return func(f *File) {
		f.strictOpen = false
	}
}

// FailOnUnknownFields tells the file decoder to fail if a key exists
// in the file but not in the destination. This feature was used in yaml.v2
// but not (yet?) developed with YAML.v3.
// Change of version was due to https://github.com/go-yaml/yaml/issues/395.
func FailOnUnknownFields() Option {
	return func(f *File) {
		f.strictUnmarshal = true
	}
}
