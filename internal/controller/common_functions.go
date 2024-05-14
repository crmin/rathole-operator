package controller

import (
	"github.com/pelletier/go-toml"
)

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) []string {
	for i, item := range slice {
		if item == s {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func ConvertSpecToToml(parent *string, spec interface{}) (string, error) {
	var (
		content []byte
		err     error
	)
	if parent == nil {
		content, err = toml.Marshal(spec)
	}
	content, err = toml.Marshal(
		map[string]interface{}{
			*parent: spec,
		},
	)
	return string(content), err
}
