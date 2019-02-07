package config

import "github.com/krostar/config/trivialerr"

type stubSourceThatUseReflection map[string][]byte

func (s stubSourceThatUseReflection) Name() string {
	return "stub"
}

func (s stubSourceThatUseReflection) GetReprValueByKey(name string) ([]byte, error) {
	for key, resp := range s {
		if name == key {
			return resp, nil
		}
	}
	return nil, trivialerr.New("not found")
}
