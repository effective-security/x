package values

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"sort"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/effective-security/xlog"
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
	return JSON(c)
}

func (c MapAny) JSONIndent() string {
	return JSONIndent(c)
}

// YAML returns YAML encoded string
func (c MapAny) YAML() string {
	raw, _ := yaml.Marshal(c)
	return string(raw)
}

// Has will return true if the map contains the key
func (c MapAny) Has(k string) bool {
	if c == nil {
		return false
	}
	_, ok := c[k]
	return ok
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

// Float64 will return the named value as a float64
func (c MapAny) Float64(k string) float64 {
	if c == nil {
		return 0
	}
	return Float64(c[k])
}

// Float32 will return the named value as a float32
func (c MapAny) Float32(k string) float32 {
	if c == nil {
		return 0
	}
	return Float32(c[k])
}

func (c MapAny) IsSlice(k string) bool {
	if c == nil {
		return false
	}
	return IsSlice(c[k])
}

func (c MapAny) Slice(k string) []any {
	if c == nil {
		return nil
	}
	return c[k].([]any)
}

func (c MapAny) IsMap(k string) bool {
	if c == nil {
		return false
	}
	return IsMap(c[k])
}

func (c MapAny) Map(k string) MapAny {
	if c == nil {
		return nil
	}
	return c[k].(map[string]any)
}

// GetOrSet returns existing value or set new value
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

// Extract extracts value from nested map using path of keys
func (c MapAny) Extract(path ...string) MapAny {
	if c == nil {
		return nil
	}

	m := c
	for _, prop := range path {
		obj, ok := m[prop]
		if !ok {
			m = nil
			break
		}
		objMap, ok := obj.(map[string]any)
		if !ok {
			m = nil
			break
		}
		m = objMap
	}

	return m
}

func (c MapAny) TraverseSubMaps(f func(k string, v MapAny) (bool, error)) error {
	if c == nil {
		return nil
	}
	for k, obj := range c {
		if IsMap(obj) {
			objMap := obj.(map[string]any)
			ok, err := f(k, objMap)
			if err != nil {
				return err
			}
			if !ok {
				return nil
			}
			subMap := MapAny(objMap)
			err = subMap.TraverseSubMaps(f)
			if err != nil {
				return err
			}
		} else if IsSlice(obj) {
			if objSlice, ok := obj.([]any); ok {
				for _, item := range objSlice {
					if IsMap(item) {
						itemMap := item.(map[string]any)
						subMap := MapAny(itemMap)
						err := subMap.TraverseSubMaps(f)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

// Merge merges maps
func (c *MapAny) Merge(m MapAny) *MapAny {
	if *c == nil {
		*c = MapAny{}
	}
	for k, v := range m {
		(*c)[k] = v
	}
	return c
}

// Scan implements the Scanner interface.
func (c *MapAny) Scan(value any) error {
	if value == nil {
		*c = nil
		return nil
	}

	var s []byte
	switch vid := value.(type) {
	case []byte:
		s = vid
	case string:
		s = []byte(vid)
	default:
		return errors.Errorf("unsupported scan type: %T", value)
	}

	if len(s) == 0 {
		*c = MapAny{}
		return nil
	}
	return errors.WithStack(json.Unmarshal(s, c))
}

// Value implements the driver Valuer interface.
func (c MapAny) Value() (driver.Value, error) {
	if len(c) == 0 {
		return nil, nil
	}
	value, err := json.Marshal(c)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return string(value), nil
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
