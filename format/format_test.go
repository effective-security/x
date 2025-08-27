package format_test

import (
	"testing"

	"github.com/effective-security/x/format"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYesNo(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "yes", format.YesNo(true))
	assert.Equal(t, "no", format.YesNo(false))
}

func TestEnabled(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "enabled", format.Enabled(true))
	assert.Equal(t, "disabled", format.Enabled(false))
}

func TestNumber(t *testing.T) {
	t.Parallel()
	require.Equal(t, "42", format.Number(42))
	require.Equal(t, "123456789", format.Number(int64(123456789)))
	require.Equal(t, "0", format.Number(uint(0)))
}

func TestFloat(t *testing.T) {
	t.Parallel()
	require.Equal(t, "3.14", format.Float(3.14159))
	require.Equal(t, "0.00", format.Float(0.0))
	require.Equal(t, "-1.23", format.Float(-1.234))
}

func TestStringMax(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "hello", format.StringMax(10, "hello"))
	assert.Equal(t, "hello", format.StringMax(5, "hello"))
	assert.Equal(t, "hello...", format.StringMax(5, "hello world"))
}

func TestStrinsgs(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "", format.Strings(nil))
	assert.Equal(t, "hello", format.Strings([]string{"hello"}))
	assert.Equal(t, "hello,world", format.Strings([]string{"hello", "world"}))
	assert.Equal(t, "hello,world", format.Strings([]string{"hello", "world", ""}))
	assert.Equal(t, "hello,world", format.Strings([]string{"hello", "", "world"}))
}

func TestStrinsgsMax(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "", format.StringsMax(3, nil))
	assert.Equal(t, "hello", format.StringsMax(3, []string{"hello"}))
	assert.Equal(t, "hello...", format.StringsMax(3, []string{"hello", ""}))
	assert.Equal(t, "hello,world", format.StringsMax(10, []string{"hello", "world"}))
	assert.Equal(t, "hello,world...", format.StringsMax(10, []string{"hello", "world", ""}))
	assert.Equal(t, "hello,world...", format.StringsMax(10, []string{"hello", "world", "foo"}))
	assert.Equal(t, "hello,world...", format.StringsMax(10, []string{"hello", "", "world", "foo"}))
	assert.Equal(t, "hello,world...", format.StringsMax(10, []string{"hello", "", "world", "", "foo"}))
}

func TestStringsAndMore(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "", format.StringsAndMore(3, nil))
	assert.Equal(t, "hello", format.StringsAndMore(3, []string{"hello"}))
	assert.Equal(t, "hello,world", format.StringsAndMore(3, []string{"hello", "world"}))
	assert.Equal(t, "hello,world", format.StringsAndMore(3, []string{"hello", "world", ""}))
	assert.Equal(t, "hello,world", format.StringsAndMore(3, []string{"hello", "", "world"}))
	assert.Equal(t, "hello,world,foo, 2 more...", format.StringsAndMore(3, []string{"hello", "world", "foo", "bar", "baz"}))
	assert.Equal(t, "hello,world,foo, 3 more...", format.StringsAndMore(3, []string{"hello", "world", "foo", "bar", "baz", ""}))
	assert.Equal(t, "hello,world, 3 more...", format.StringsAndMore(3, []string{"hello", "", "world", "foo", "bar", "baz"}))
	assert.Equal(t, "hello,world, 4 more...", format.StringsAndMore(3, []string{"hello", "", "world", "", "foo", "bar", "baz"}))
}

func Test_DisplayName(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "Test", format.DisplayName("Test"))
	assert.Equal(t, "test", format.DisplayName("test"))
	assert.Equal(t, "Test Data", format.DisplayName("TestData"))
	assert.Equal(t, "AWS Name", format.DisplayName("AWSName"))
	assert.Equal(t, "S3 Location", format.DisplayName("S3Location"))
	assert.Equal(t, "EC2 Instance", format.DisplayName("EC2Instance"))
	assert.Equal(t, "Asset ID", format.DisplayName("AssetID"))
	assert.Equal(t, "IDs", format.DisplayName("IDs"))
	assert.Equal(t, "Asset IDs", format.DisplayName("AssetIDs"))
	assert.Equal(t, "API Key", format.DisplayName("APIKey"))
	assert.Equal(t, "API Keys", format.DisplayName("APIKeys"))
	assert.Equal(t, "GUID", format.DisplayName("GUID"))
	assert.Equal(t, "GUIDs", format.DisplayName("GUIDs"))
}

func Test_DisplayName_AcronymPlurals(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   string
		want string
	}{
		{in: "IDs", want: "IDs"},
		{in: "AssetIDs", want: "Asset IDs"},
		{in: "UserIDs", want: "User IDs"},
		{in: "UserIDsCount", want: "User IDs Count"},
		{in: "APIs", want: "APIs"},
		{in: "URLs", want: "URLs"},
		{in: "GUIDs", want: "GUIDs"},
		{in: "AssetGUIDs", want: "Asset GUIDs"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, format.DisplayName(tt.in))
		})
	}
}

func Test_TextOneLine(t *testing.T) {
	t.Parallel()
	doc := `line 1
with continuation
	
Line 2.

Line 3
`
	assert.Equal(t, `line 1 with continuation. Line 2. Line 3.`, format.TextOneLine(doc))

	exp := `  line 1
  with continuation
  	
  Line 2.
  
  Line 3
  
`
	assert.Equal(t, exp, format.TextWithIndent(doc, "  ", false))

	exp = `line 1
  with continuation
  	
  Line 2.
  
  Line 3
  
`
	assert.Equal(t, exp, format.TextWithIndent(doc, "  ", true))
}
