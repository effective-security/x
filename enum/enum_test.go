package enum_test

import (
	"strings"
	"testing"

	"github.com/effective-security/x/enum"
	"github.com/stretchr/testify/assert"
)

type Severity_Enum int32

const (
	Severity_Invalid  Severity_Enum = 0
	Severity_Unknown  Severity_Enum = 1
	Severity_Low      Severity_Enum = 2
	Severity_Medium   Severity_Enum = 4
	Severity_High     Severity_Enum = 8
	Severity_Critical Severity_Enum = 16
	Severity_All      Severity_Enum = 0x7fffffff
)

// Enum value maps for Severity_Enum.
var (
	Severity_Enum_name = map[int32]string{
		0:          "Invalid",
		1:          "Unknown",
		2:          "Low",
		4:          "Medium",
		8:          "High",
		16:         "Critical",
		0x7fffffff: "All",
	}
	Severity_Enum_value = map[string]int32{
		"Invalid":  0,
		"Unknown":  1,
		"Low":      2,
		"Medium":   4,
		"High":     8,
		"Critical": 16,
		"All":      0x7fffffff,
	}
)

func (s Severity_Enum) ValuesMap() map[string]int32 {
	return Severity_Enum_value
}
func (s Severity_Enum) NamesMap() map[int32]string {
	return Severity_Enum_name
}
func (s Severity_Enum) SupportedNames() string {
	return enum.SupportedNames[Severity_Enum]()
}
func (s Severity_Enum) ValueNames() []string {
	return enum.FlagNames(s)
}
func (s Severity_Enum) ValueString() string {
	return strings.Join(s.ValueNames(), ",")
}
func (s Severity_Enum) Flags() []Severity_Enum {
	return enum.Flags(s)
}
func (s Severity_Enum) FlagsInt() []int32 {
	return enum.FlagsInt(s)
}

func Test_SupportedNames(t *testing.T) {
	assert.Equal(t, "All,Critical,High,Invalid,Low,Medium,Unknown", Severity_All.SupportedNames())

	assert.Equal(t, Severity_Critical|Severity_High, enum.Parse[Severity_Enum]("Critical|High"))
	assert.Equal(t, Severity_Critical, enum.Parse[Severity_Enum]("Critical"))
	assert.Equal(t, Severity_Critical|Severity_High, enum.Parse[Severity_Enum]("Critical,High"))
	assert.Equal(t, Severity_Critical|Severity_High, enum.Convert[Severity_Enum]([]string{"Critical", "High"}))

	e := Severity_Critical | Severity_High
	assert.Equal(t, "High,Critical", e.ValueString())
	assert.Equal(t, []string{"High", "Critical"}, e.ValueNames())
	assert.Equal(t, []Severity_Enum{Severity_High, Severity_Critical}, e.Flags())
	assert.Equal(t, []int32{8, 16}, e.FlagsInt())

	e = Severity_All
	assert.Equal(t, "Unknown,Low,Medium,High,Critical", e.ValueString())
	assert.Equal(t, []string{"Unknown", "Low", "Medium", "High", "Critical"}, e.ValueNames())
	assert.Equal(t, []Severity_Enum{Severity_Unknown, Severity_Low, Severity_Medium, Severity_High, Severity_Critical}, e.Flags())
	assert.Equal(t, []int32{1, 2, 4, 8, 16}, e.FlagsInt())
}
