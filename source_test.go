package config

import (
	"github.com/krostar/config/trivialerr"
	"github.com/pkg/errors"
)

type stubSourceThatUseReflection map[string][]byte

func (s stubSourceThatUseReflection) Name() string {
	return "stub reflect"
}

func (s stubSourceThatUseReflection) GetReprValueByKey(name string) ([]byte, error) {
	for key, resp := range s {
		if name == key {
			return resp, nil
		}
	}
	return nil, trivialerr.New("not found")
}

type stubSourceThatUnmarshal int

func (s stubSourceThatUnmarshal) Name() string {
	return "stub unmarshal"
}

func (s stubSourceThatUnmarshal) Unmarshal(interface{}) error {
	var err error

	switch {
	case s < 0:
		err = trivialerr.New("not found")
	case s > 0:
		err = errors.New("not found")
	}
	return err
}
