package values

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/effective-security/xlog"
)

var logger = xlog.NewPackageLogger("github.com/effective-security/x", "values")

// String returns string value
func String(v any) string {
	if v == nil {
		return ""
	}
	switch tv := v.(type) {
	case string:
		return tv
	case []string:
		if len(tv) > 0 {
			return tv[0]
		}
		return ""
	case []any:
		if len(tv) > 0 {
			return String(v.([]any)[0])
		}
		return ""
	case int:
		return strconv.Itoa(tv)
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
		return strconv.FormatUint(uint64(tv), 10)
	case float64:
		return strconv.FormatUint(uint64(tv), 10)
	case bool:
		return Select(tv, "true", "false")
	default:
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

// IsCollection returns true for slices and maps
func IsCollection(value any) bool {
	if value == nil {
		return false
	}
	kind := reflect.TypeOf(value).Kind()
	return kind == reflect.Slice || kind == reflect.Map
}

// JSON returns the value as a JSON string
func JSON(value any) string {
	if value == nil {
		return ""
	}
	b, _ := json.Marshal(value)
	return string(b)
}
