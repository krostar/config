package configue

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type customDuration time.Duration

// ToDuration converts the custom duration back to the real time.Duration.
func (cd *customDuration) ToDuration() time.Duration {
	return time.Duration(*cd)
}

// ToInt64 converts the custom curation to the original value behind a time.Duration.
func (cd *customDuration) ToInt64() int64 {
	return int64(*cd)
}

// MarshalJSON implements json Marshaler interface.
func (cd *customDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal((time.Duration(*cd)).String())
}

// UnMarshalJSON implements json Unmarshaler interface.
func (cd *customDuration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return errors.Wrap(err, "json failed to unmarshal to v")
	}

	var (
		d   time.Duration
		err error
	)

	switch value := v.(type) {
	case float64:
		d = time.Duration(value)
	case string:
		d, err = time.ParseDuration(value)
	default:
		err = errors.New("invalid duration type")
	}

	if err != nil {
		return err
	}

	*cd = customDuration(d)
	return nil
}
