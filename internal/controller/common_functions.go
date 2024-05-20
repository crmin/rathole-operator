package controller

import (
	"bytes"
	"github.com/pelletier/go-toml/v2"
	"math/rand"
	"time"
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
		//content, err = toml.Marshal(spec)
		content, err = TomlMarshal(spec)
	}

	//content, err = toml.Marshal(
	//	map[string]interface{}{
	//		*parent: spec,
	//	},
	//)
	content, err = TomlMarshal(
		map[string]interface{}{
			*parent: spec,
		},
	)
	return string(content), err
}

func TomlMarshal(data interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := toml.NewEncoder(&buf)
	enc.SetIndentTables(false)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetSuffix(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"

	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}

	return string(result)
}
