package values

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/effective-security/x/enum"
	"github.com/effective-security/xlog"
)

var logger = xlog.NewPackageLogger("github.com/effective-security/x", "values")

type HasDisplayName interface {
	DisplayName() string
}

type HasDisplayNames interface {
	DisplayNames() []string
}

type HasName interface {
	Name() string
}

type HasNames interface {
	Names() []string
}

type HasValuesMap interface {
	ValuesMap() map[string]int32
}

type HasNamesMap interface {
	NamesMap() map[int32]string
}

type HasDisplayNamesMap interface {
	DisplayNamesMap() map[int32]string
}

// String returns string value
func String(v any) string {
	if v == nil {
		return ""
	}
	switch tv := v.(type) {
	case string:
		return tv
	case []string:
		return strings.Join(tv, ",")
	case []any:
		list := make([]string, len(tv))
		for i, v := range tv {
			list[i] = String(v)
		}
		return strings.Join(list, ",")
	case int:
		return strconv.Itoa(tv)
	case int16:
		return strconv.Itoa(int(tv))
	case int32:
		return strconv.Itoa(int(tv))
	case int64:
		return strconv.FormatInt(tv, 10)
	case uint:
		return strconv.FormatUint(uint64(tv), 10)
	case uint32:
		return strconv.FormatUint(uint64(tv), 10)
	case uint64:
		return strconv.FormatUint(uint64(tv), 10)
	case float32:
		return strconv.FormatFloat(float64(tv), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(tv, 'f', -1, 64)
	case bool:
		return Select(tv, "true", "false")
	case HasDisplayName:
		return tv.DisplayName()
	case HasName:
		return tv.Name()
	case fmt.Stringer:
		return tv.String()
	default:
		kind := reflect.TypeOf(v).Kind()
		if kind == reflect.Slice {
			v := reflect.ValueOf(v)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			list := make([]string, v.Len())
			for i := 0; i < v.Len(); i++ {
				list[i] = String(v.Index(i).Interface())
			}
			return strings.Join(list, ",")
		}
		logger.KV(xlog.DEBUG, "reason", "unsupported", "type", fmt.Sprintf("%T", v))
		return xlog.EscapedString(v)
	}
}

// StringSlice returns slice of string values
func StringSlice(v any) []string {
	if v == nil {
		return []string{}
	}

	switch tv := v.(type) {
	case string:
		return strings.Split(tv, ",")
	case []string:
		return tv
	case []any:
		list := make([]string, len(tv))
		for i, v := range tv {
			list[i] = String(v)
		}
		return list
	default:
		logger.KV(xlog.DEBUG, "reason", "unsupported", "type", reflect.TypeOf(v))
		return []string{}
	}
}

// IntSlice returns slice of int values
func IntSlice(v any) []int {
	if v == nil {
		return []int{}
	}

	switch tv := v.(type) {
	case []int:
		return tv
	case []int32:
		list := make([]int, len(tv))
		for i, v := range tv {
			list[i] = int(v)
		}
		return list
	case []int64:
		list := make([]int, len(tv))
		for i, v := range tv {
			list[i] = int(v)
		}
		return list
	case []uint:
		list := make([]int, len(tv))
		for i, v := range tv {
			list[i] = int(v)
		}
		return list
	case []uint32:
		list := make([]int, len(tv))
		for i, v := range tv {
			list[i] = int(v)
		}
		return list
	case []uint64:
		list := make([]int, len(tv))
		for i, v := range tv {
			list[i] = int(v)
		}
		return list
	case []float32:
		list := make([]int, len(tv))
		for i, v := range tv {
			list[i] = int(v)
		}
		return list
	case []float64:
		list := make([]int, len(tv))
		for i, v := range tv {
			list[i] = int(v)
		}
		return list
	case string:
		return IntSlice(strings.Split(tv, ","))
	case []string:
		list := make([]int, len(tv))
		for i, v := range tv {
			list[i] = Int(v)
		}
		return list
	case []any:
		list := make([]int, len(tv))
		for i, v := range tv {
			list[i] = Int(v)
		}
		return list
	default:
		kind := reflect.TypeOf(v).Kind()
		if kind == reflect.Slice {
			v := reflect.ValueOf(v)
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			list := make([]int, v.Len())
			for i := 0; i < v.Len(); i++ {
				list[i] = Int(v.Index(i).Interface())
			}
			return list
		}

		logger.KV(xlog.DEBUG, "reason", "unsupported", "type", reflect.TypeOf(v))
		return []int{}
	}
}

// Bool returns bool value
func Bool(v any) bool {
	if v == nil {
		return false
	}
	switch tv := v.(type) {
	case bool:
		return tv
	case string:
		return tv == "true" || tv == "yes"
	case []any:
		if len(tv) > 0 {
			return Bool(v.([]any)[0])
		}
		return false
	default:
		logger.KV(xlog.DEBUG, "reason", "unsupported", "type", fmt.Sprintf("%T", v))
		return false
	}
}

// Time will return the value as Time
func Time(v any) *time.Time {
	if v == nil {
		return nil
	}
	switch tv := v.(type) {
	case []any:
		if len(tv) > 0 {
			return Time(v.([]any)[0])
		}
		return nil
	case time.Time:
		return &tv
	case *time.Time:
		return tv
	case int64:
		t := time.Unix(tv, 0)
		return &t
	case uint64:
		t := time.Unix(int64(tv), 0)
		return &t
	case float64:
		t := time.Unix(int64(tv), 0)
		return &t
	case int:
		t := time.Unix(int64(tv), 0)
		return &t
	case string:
		if len(tv) > 20 {
			t, err := time.Parse("2006-01-02T15:04:05.000-0700", tv)
			if err != nil {
				return nil
			}
			return &t
		}
		unix, err := strconv.ParseInt(tv, 10, 64)
		if err != nil {
			logger.KV(xlog.DEBUG, "val", v, "err", err.Error())
			return nil
		}
		t := time.Unix(unix, 0)
		return &t
	default:
		logger.KV(xlog.DEBUG, "reason", "unsupported", "val", v, "type", fmt.Sprintf("%T", v))
		return nil
	}
}

// Int will return the value as int
func Int(v any) int {
	if v == nil {
		return 0
	}
	switch tv := v.(type) {
	case []any:
		if len(tv) > 0 {
			return Int(v.([]any)[0])
		}
		return 0
	case int:
		return tv
	case int16:
		return int(tv)
	case int32:
		return int(tv)
	case int64:
		return int(tv)
	case uint:
		return int(tv)
	case uint32:
		return int(tv)
	case uint64:
		return int(tv)
	case float32:
		return int(tv)
	case float64:
		return int(tv)
	case string:
		i, err := strconv.Atoi(tv)
		if err != nil {
			logger.KV(xlog.DEBUG, "val", v, "err", err.Error())
			return 0
		}
		return i
	case enum.ProtoEnum:
		return int(tv.Number())
	default:
		logger.KV(xlog.DEBUG, "reason", "unsupported", "val", v, "type", fmt.Sprintf("%T", v))
		return 0
	}
}

// UInt64 will return the value as uint64
func UInt64(v any) uint64 {
	if v == nil {
		return 0
	}
	switch tv := v.(type) {
	case []any:
		if len(tv) > 0 {
			return UInt64(v.([]any)[0])
		}
		return 0
	case int:
		return uint64(tv)
	case int16:
		return uint64(tv)
	case int32:
		return uint64(tv)
	case int64:
		return uint64(tv)
	case uint:
		return uint64(tv)
	case uint32:
		return uint64(tv)
	case uint64:
		return uint64(tv)
	case float32:
		return uint64(tv)
	case float64:
		return uint64(tv)
	case string:
		i64, err := strconv.ParseUint(tv, 10, 64)
		if err != nil {
			logger.KV(xlog.DEBUG, "val", v, "err", err.Error())
			return 0
		}
		return i64
	default:
		logger.KV(xlog.DEBUG, "reason", "unsupported", "val", v, "type", fmt.Sprintf("%T", v))
		return 0
	}
}

// Int64 will return the value as int64
func Int64(v any) int64 {
	if v == nil {
		return 0
	}
	switch tv := v.(type) {
	case []any:
		if len(tv) > 0 {
			return Int64(v.([]any)[0])
		}
		return 0
	case int:
		return int64(tv)
	case int16:
		return int64(tv)
	case int32:
		return int64(tv)
	case int64:
		return int64(tv)
	case uint:
		return int64(tv)
	case uint32:
		return int64(tv)
	case uint64:
		return int64(tv)
	case float32:
		return int64(tv)
	case float64:
		return int64(tv)
	case string:
		i64, err := strconv.ParseInt(tv, 10, 64)
		if err != nil {
			logger.KV(xlog.DEBUG, "val", v, "err", err.Error())
			return 0
		}
		return i64
	default:
		logger.KV(xlog.DEBUG, "reason", "unsupported", "val", v, "type", fmt.Sprintf("%T", v))
		return 0
	}
}

// Float32 will return the value as float32
func Float32(v any) float32 {
	if v == nil {
		return 0
	}
	switch tv := v.(type) {
	case []any:
		if len(tv) > 0 {
			return Float32(v.([]any)[0])
		}
		return 0
	case int:
		return float32(tv)
	case int16:
		return float32(tv)
	case int32:
		return float32(tv)
	case int64:
		return float32(tv)
	case uint:
		return float32(tv)
	case uint32:
		return float32(tv)
	case uint64:
		return float32(tv)
	case float32:
		return tv
	case float64:
		return float32(tv)
	case string:
		i64, err := strconv.ParseFloat(tv, 32)
		if err != nil {
			logger.KV(xlog.DEBUG, "val", v, "err", err.Error())
			return 0
		}
		return float32(i64)
	default:
		logger.KV(xlog.DEBUG, "reason", "unsupported", "val", v, "type", fmt.Sprintf("%T", v))
		return 0
	}
}

// Float64 will return the value as float64
func Float64(v any) float64 {
	if v == nil {
		return 0
	}
	switch tv := v.(type) {
	case []any:
		if len(tv) > 0 {
			return Float64(v.([]any)[0])
		}
		return 0
	case int:
		return float64(tv)
	case int16:
		return float64(tv)
	case int32:
		return float64(tv)
	case int64:
		return float64(tv)
	case uint:
		return float64(tv)
	case uint32:
		return float64(tv)
	case uint64:
		return float64(tv)
	case float32:
		return float64(tv)
	case float64:
		return tv
	case string:
		i64, err := strconv.ParseFloat(tv, 64)
		if err != nil {
			logger.KV(xlog.DEBUG, "val", v, "err", err.Error())
			return 0
		}
		return float64(i64)
	default:
		logger.KV(xlog.DEBUG, "reason", "unsupported", "val", v, "type", fmt.Sprintf("%T", v))
		return 0
	}
}

// IsCollection returns true for slices and maps
func IsCollection(value any) bool {
	if value == nil {
		return false
	}
	kind := reflect.TypeOf(value).Kind()
	return kind == reflect.Slice || kind == reflect.Map
}

func IsSlice(value any) bool {
	if value == nil {
		return false
	}
	kind := reflect.TypeOf(value).Kind()
	return kind == reflect.Slice
}

func IsMap(value any) bool {
	if value == nil {
		return false
	}
	kind := reflect.TypeOf(value).Kind()
	return kind == reflect.Map
}

// JSON returns the value as a JSON string
func JSON(value any) string {
	if value == nil {
		return ""
	}
	b, _ := json.Marshal(value)
	return string(b)
}

// JSONIndent returns the value as a JSON string with indentation
func JSONIndent(value any) string {
	if value == nil {
		return ""
	}
	b, _ := json.MarshalIndent(value, "", "\t")
	return string(b)
}

// IndentJSON indents a JSON string
// It will return the original string if it fails to indent
// This is useful for pretty printing JSON strings
func IndentJSON(data string) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(data), "", "\t"); err == nil {
		return buf.String()
	}
	return data
}

// IsEmpty checks if a value is considered empty
func IsEmpty(value any) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	if v.IsZero() {
		return true
	}
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	case reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return IsEmpty(v.Elem().Interface())
	case reflect.Interface:
		elem := v.Elem()
		if !elem.IsValid() {
			return true
		}
		return IsEmpty(elem.Interface())
	default:
		// For other types (structs, channels, etc.), consider them non-empty
		return false
	}
}

// Shrink removes empty values recursively:
// - returns nil if value is empty
// - if it's a Map, then only non empty values are returned in a map, or nil
// - if it's a Slice, then only non empty values are returned.
func Shrink(value any) any {
	if IsEmpty(value) {
		return nil
	}

	v := reflect.ValueOf(value)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		return Shrink(v.Elem().Interface())
	}

	switch v.Kind() {
	case reflect.Map:
		// Handle map[string]any specifically for MapAny compatibility
		if mapVal, ok := value.(map[string]any); ok {
			shrunk := MapAny(mapVal).Shrink()
			if shrunk == nil {
				return nil
			}
			// Convert back to map[string]any to maintain type consistency
			return map[string]any(shrunk)
		}

		// Handle MapAny type directly
		if mapAny, ok := value.(MapAny); ok {
			shrunk := mapAny.Shrink()
			if shrunk == nil {
				return nil
			}
			// Convert back to map[string]any to maintain type consistency
			return map[string]any(shrunk)
		}

		// Handle other map types generically
		newMap := reflect.MakeMap(v.Type())
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			if !val.IsValid() {
				continue
			}
			shrunkVal := Shrink(val.Interface())
			if shrunkVal != nil {
				newMap.SetMapIndex(key, reflect.ValueOf(shrunkVal))
			}
		}
		if newMap.Len() == 0 {
			return nil
		}
		return newMap.Interface()

	case reflect.Slice, reflect.Array:
		var newSlice []any
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			if !elem.IsValid() {
				continue
			}
			shrunkElem := Shrink(elem.Interface())
			if shrunkElem != nil {
				newSlice = append(newSlice, shrunkElem)
			}
		}
		if len(newSlice) == 0 {
			return nil
		}
		return newSlice

	default:
		// For primitive types and other non-collection types,
		// return the value as-is if it's not empty
		return value
	}
}
