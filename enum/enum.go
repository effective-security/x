package enum

import (
	"sort"
	"strings"
)

const nameSeparator = ","

// Enum interface for generic enum
type Enum interface {
	~int32 | ~uint32
}

// Names interface for enum with names
type Names interface {
	Enum
	NamesMap() map[int32]string
}

// Values interface for enum with values
type Values interface {
	Enum
	ValuesMap() map[string]int32
}

// SupportedNames returns supported Enum values concatenated by ","
func SupportedNames[E Values]() string {
	var e E
	return NamesHelpString(e.ValuesMap())
}

// NamesHelpString returns supported Enum values concatenated by ","
func NamesHelpString(vals map[string]int32) string {
	var typs []string
	for typ := range vals {
		typs = append(typs, typ)
	}
	sort.Strings(typs)
	return strings.Join(typs, nameSeparator)
}

// Convert returns enum value from names
func Convert[E Values](names []string) E {
	var res E
	values := res.ValuesMap()
	for _, name := range names {
		res |= E(values[name])
	}
	return res
}

// Parse returns enum value from names
func Parse[E Values](val string) E {
	var res E
	values := res.ValuesMap()

	var names []string
	if strings.Contains(val, "|") {
		names = strings.Split(val, "|")
	} else if strings.Contains(val, nameSeparator) {
		names = strings.Split(val, nameSeparator)
	} else {
		names = []string{val}
	}

	for _, name := range names {
		res |= E(values[name])
	}
	return res
}

// FlagNames returns list of enum value names from flag value
func FlagNames[E Names](val E) []string {
	names := val.NamesMap()

	var vals []string
	for i := E(1); i <= val; i <<= 1 {
		if val&i == i {
			name := names[int32(i)]
			if name != "" {
				vals = append(vals, names[int32(i)])
			}
		}
	}
	return vals
}

// FlagsInt returns list of enum values from flag
func FlagsInt[E Names](val E) []int32 {
	var vals []int32
	for i := E(1); i <= val; i <<= 1 {
		if val&i == i {
			vals = append(vals, int32(i))
		}
	}
	return vals
}

// FlagsInt returns list of enum values from flag
func Flags[E Names](val E) []E {
	var vals []E
	for i := E(1); i <= val; i <<= 1 {
		if val&i == i {
			vals = append(vals, i)
		}
	}
	return vals
}
