package configloader

import (
	"os"
	"reflect"
	"strings"

	"github.com/cockroachdb/errors"
)

// Expander is used to expand variables in the input object
type Expander struct {
	Variables      map[string]string
	SecretProvider SecretProvider
}

// ExpandAll replace variables in the input object, using default Expander.
// The input object must be a pointer to a struct.
// If secrets are used, SecretProviderInstance must be set.
// Values that start with env://, file:// or secret:// are resolved.
// Names inside ${} are looked up in Variables, then os.Getenv; a missing
// name becomes an empty string, matching os.Expand. After that pass, any
// remaining "${" is an error (typically a value that expanded to another
// interpolation). A failed or missing SecretProvider for ${secret://…}
// is always an error.
func ExpandAll(obj any) error {
	e := Expander{SecretProvider: SecretProviderInstance}
	return e.ExpandAll(obj)
}

// ExpandAll replace variables in the input object
func (f *Expander) ExpandAll(obj any) error {
	return f.doSubstituteEnvVars(reflect.ValueOf(obj))
}

// Expand replace variables in the input string
func (f *Expander) Expand(s string) (string, error) {
	var expandErr error
	if strings.Contains(s, "${") {
		s = os.Expand(s, func(env string) string {
			if strings.HasPrefix(env, SecretSource) {
				if f.SecretProvider == nil {
					if expandErr == nil {
						expandErr = errors.Errorf("secret loader not provided: unable to expand: ${%s}", env)
					}
					return ""
				}
				name := strings.TrimPrefix(env, SecretSource)
				sec, err := f.SecretProvider.GetSecret(name)
				if err != nil {
					if expandErr == nil {
						expandErr = errors.WithMessagef(err, "unable to load secret: %s", name)
					}
					return ""
				}
				return sec
			}

			if va, ok := f.Variables[env]; ok {
				return va
			}
			return os.Getenv(env)
		})
		if expandErr != nil {
			return s, expandErr
		}
	}

	if strings.Contains(s, "${") {
		return s, errors.Errorf("unable to resolve variables: %s", s)
	}

	s, err := ResolveValueWithSecrets(s, f.SecretProvider)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (f *Expander) doSubstituteEnvVars(v reflect.Value) error {
	if v.Kind() == reflect.Pointer {
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
	case reflect.Pointer:
		if err := f.doSubstituteEnvVars(v.Elem()); err != nil {
			return err
		}
	case reflect.Map:
		if v.IsNil() {
			return nil
		}
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			if !val.IsValid() {
				continue
			}
			if err := f.expandMapEntry(v, key, val); err != nil {
				return err
			}
		}
	default:
	}
	return nil
}

// expandMapEntry expands a map value and writes it back. MapRange / MapIndex
// values are not addressable, so strings and structs must be copied.
func (f *Expander) expandMapEntry(m, key, val reflect.Value) error {
	switch val.Kind() {
	case reflect.Interface:
		if val.IsNil() {
			return nil
		}
		return f.expandMapEntry(m, key, val.Elem())
	case reflect.String:
		expanded, err := f.Expand(val.String())
		if err != nil {
			return err
		}
		newVal := reflect.ValueOf(expanded)
		if newVal.Type() != val.Type() && newVal.CanConvert(val.Type()) {
			newVal = newVal.Convert(val.Type())
		}
		m.SetMapIndex(key, newVal)
		return nil
	case reflect.Pointer:
		return f.doSubstituteEnvVars(val)
	default:
		if val.CanAddr() {
			return f.doSubstituteEnvVars(val)
		}
		cpy := reflect.New(val.Type()).Elem()
		cpy.Set(val)
		if err := f.doSubstituteEnvVars(cpy); err != nil {
			return err
		}
		m.SetMapIndex(key, cpy)
		return nil
	}
}
