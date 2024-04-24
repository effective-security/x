package values

import (
	"bytes"
	"encoding/json"
	"sort"
	"time"

	"github.com/effective-security/xlog"
	"github.com/pkg/errors"
	"golang.org/x/exp/constraints"
	"gopkg.in/yaml.v3"
)

// MapAny provides map of values
type MapAny map[string]any

// FromJSON returns map from json encoded string,
// this method does not return value on error, as the value is expected a valid JSON string.
func FromJSON(s string) MapAny {
	var val MapAny
	if s != "" && s != "{}" && s != "[]" {
		err := json.Unmarshal([]byte(s), &val)
		if err != nil {
			logger.KV(xlog.DEBUG,
				"reason", "unmarshal",
				"val", s,
				"err", err.Error())
		}
	}
	return val
}

// FromYAML returns map from yaml encoded string,
// this method does not return value on error, as the value is expected a valid YAML string.
func FromYAML(s string) MapAny {
	var val MapAny
	if s != "" && s != "{}" && s != "[]" {
		err := yaml.Unmarshal([]byte(s), &val)
		if err != nil {
			logger.KV(xlog.DEBUG,
				"reason", "unmarshal",
				"val", s,
				"err", err.Error())
		}
	}
	return val
}

// To converts the values to the value pointed to by val.
func (c MapAny) To(val any) error {
	raw, err := json.Marshal(c)
	if err != nil {
		return errors.WithStack(err)
	}

	d := json.NewDecoder(bytes.NewReader(raw))
	if err := d.Decode(val); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// JSON returns JSON encoded string
func (c MapAny) JSON() string {
	raw, _ := json.Marshal(c)
	return string(raw)
}

// YAML returns YAML encoded string
func (c MapAny) YAML() string {
	raw, _ := yaml.Marshal(c)
	return string(raw)
}

// String will return the value as a string,
// if the underlying type is not a string,
// it will try and co-oerce it to a string.
func (c MapAny) String(k string) string {
	if c == nil {
		return ""
	}
	return String(c[k])
}

// StringSlice returns slice of string values
func (c MapAny) StringSlice(k string) []string {
	if c == nil {
		return nil
	}
	return StringSlice(c[k])
}

// Bool will return the value as Bool
func (c MapAny) Bool(k string) bool {
	if c == nil {
		return false
	}
	return Bool(c[k])
}

// Time will return the value as Time
func (c MapAny) Time(k string) *time.Time {
	if c == nil {
		return nil
	}
	return Time(c[k])
}

// Int will return the value as an int
func (c MapAny) Int(k string) int {
	if c == nil {
		return 0
	}
	return Int(c[k])
}

// UInt64 will return the named value as an uint64
func (c MapAny) UInt64(k string) uint64 {
	if c == nil {
		return 0
	}
	return UInt64(c[k])
}

// Int64 will return the named value as an int64
func (c MapAny) Int64(k string) int64 {
	if c == nil {
		return 0
	}
	return Int64(c[k])
}

// GetOrSet existing existing value or set new value
func (c MapAny) GetOrSet(key string, getter func(key string) any) any {
	if c == nil {
		return nil
	}
	if v, ok := c[key]; ok {
		return v
	}
	v := getter(key)
	c[key] = v
	return v
}

// OrderedMapKeys returns ordered keys
func OrderedMapKeys[K constraints.Ordered, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	sort.Slice(r, func(i, j int) bool {
		return r[i] < r[j]
	})
	return r
}

// RangeOrderedMap range over ordered map
func RangeOrderedMap[K constraints.Ordered, V any](c map[K]V, f func(k K, v V) bool) {
	if c == nil {
		return
	}

	for _, k := range OrderedMapKeys(c) {
		if !f(k, c[k]) {
			break
		}
	}
}
