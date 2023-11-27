package configloader

import (
	"os"
	"reflect"
	"strings"

	"github.com/effective-security/xlog"
	"github.com/pkg/errors"
)

// Expander is used to expand variables in the input object
type Expander struct {
	Variables      map[string]string
	SecretProvider SecretProvider
}

// ExpandAll replace variables in the input object, using default Expander.
// The input object must be a pointer to a struct.
// If secrets are used, SecretProviderInstance must be set.
// The values started with env:// , file:// or secret:// must be resolved.
// The values inside ${} will be tried to be resolved,
// if not found will be substiduted with empy values as per os.Getenv function.
func ExpandAll(obj interface{}) error {
	e := Expander{SecretProvider: SecretProviderInstance}
	return e.ExpandAll(obj)
}

// ExpandAll replace variables in the input object
func (f *Expander) ExpandAll(obj interface{}) error {
	return f.doSubstituteEnvVars(reflect.ValueOf(obj))
}

// Expand replace variables in the input string
func (f *Expander) Expand(s string) (string, error) {
	// try first prefix
	s, err := ResolveValueWithSecrets(s, f.SecretProvider)
	if err != nil {
		return s, err
	}

	if strings.Contains(s, "${") {
		for key, value := range f.Variables {
			s = strings.Replace(s, key, value, -1)
		}
	}

	if strings.Contains(s, "${") {
		s = os.Expand(s, func(env string) string {
			if strings.HasPrefix(env, SecretSource) && f.SecretProvider != nil {
				name := strings.TrimPrefix(env, SecretSource)
				sec, err := f.SecretProvider.GetSecret(name)
				if err != nil {
					logger.KV(xlog.ERROR, "secret", name, "err", err.Error())
				}
				return sec
			}

			if va, ok := f.Variables[env]; ok {
				return va
			}
			return os.Getenv(env)
		})
	}
	if strings.Contains(s, "${") {
		return s, errors.Errorf("unable to resolve variables: %s", s)
	}
	return s, nil
}

func (f *Expander) doSubstituteEnvVars(v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if err := f.doSubstituteEnvVars(v.Field(i)); err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if err := f.doSubstituteEnvVars(v.Index(i)); err != nil {
				return err
			}
		}
	case reflect.String:
		if v.CanSet() {
			val, err := f.Expand(v.String())
			if err != nil {
				return err
			}
			v.SetString(val)
		}
	case reflect.Ptr:
		if err := f.doSubstituteEnvVars(v.Elem()); err != nil {
			return err
		}
	case reflect.Map:
		if v.Type().String() == "map[string]string" {
			m := v.Interface().(map[string]string)
			for k, v := range m {
				val, err := f.Expand(v)
				if err != nil {
					return err
				}
				m[k] = val
			}
		} else {
			iter := v.MapRange()
			for iter.Next() {
				if err := f.doSubstituteEnvVars(iter.Value()); err != nil {
					return err
				}
			}
		}
	default:
	}
	return nil
}
