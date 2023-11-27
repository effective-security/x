package configloader

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	// FileSource specifies to load config from a file
	FileSource = "file://"
	// EnvSource specifies to load config from an environment variable
	EnvSource = "env://"
	// SecretSource specifies to load config from a secret manager
	SecretSource = "secret://"
)

// SecretProvider is an interface to provide secrets
type SecretProvider interface {
	GetSecret(name string) (string, error)
}

// SecretProviderInstance is a global instance of SecretLoader
var SecretProviderInstance SecretProvider

// ResolveValue returns value loaded from file:// or env://
// If val does not start with file:// or env://, then the value is returned as is
func ResolveValue(val string) (string, error) {
	return ResolveValueWithSecrets(val, SecretProviderInstance)
}

// ResolveValue returns value loaded from file:// or env://
// If val does not start with file:// or env://, then the value is returned as is
func ResolveValueWithSecrets(val string, loader SecretProvider) (string, error) {
	if strings.HasPrefix(val, FileSource) {
		fn := strings.TrimPrefix(val, FileSource)
		f, err := os.ReadFile(fn)
		if err != nil {
			return val, errors.WithStack(err)
		}
		// file content
		val = string(f)
	} else if strings.HasPrefix(val, EnvSource) {
		env := strings.TrimPrefix(val, EnvSource)
		// ENV content
		val = os.Getenv(env)
		if val == "" {
			return "", errors.Errorf("environment variable not set: %s", env)
		}
	} else if strings.HasPrefix(val, SecretSource) {
		if loader == nil {
			return "", errors.Errorf("secret loader not provided")
		}
		name := strings.TrimPrefix(val, SecretSource)
		sec, err := loader.GetSecret(name)
		if err != nil {
			return val, errors.WithMessage(err, "unable to load secret")
		}
		val = sec
	}

	return val, nil
}

// UnmarshalAndExpand load JSON or YAML file to an interface and expands variables
func UnmarshalAndExpand(file string, v interface{}) error {
	err := Unmarshal(file, v)
	if err != nil {
		return err
	}

	return ExpandAll(v)
}

// Unmarshal JSON or YAML file to an interface
func Unmarshal(file string, v interface{}) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return errors.WithMessagef(err, "unable to read file")
	}

	if strings.HasSuffix(file, ".json") {
		err = json.Unmarshal(b, v)
		if err != nil {
			return errors.WithMessagef(err, "unable parse JSON: %s", file)
		}
	} else {
		err = yaml.Unmarshal(b, v)
		if err != nil {
			return errors.WithMessagef(err, "unable parse YAML: %s", file)
		}
	}
	return nil
}

// Marshal saves object to file
func Marshal(fn string, value interface{}) error {
	var data []byte
	var err error
	if strings.HasSuffix(fn, ".json") {
		data, err = json.MarshalIndent(value, "", "  ")
	} else {
		data, err = yaml.Marshal(value)
	}

	if err != nil {
		return errors.WithMessage(err, "failed to encode")
	}

	return os.WriteFile(fn, data, os.ModePerm)
}
