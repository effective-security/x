package print_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/effective-security/x/print"
	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	var buf bytes.Buffer
	print.JSON(&buf, map[string]string{"key": "value"})
	assert.JSONEq(t, `{"key": "value"}`, buf.String())
}

func TestYaml(t *testing.T) {
	var buf bytes.Buffer
	print.Yaml(&buf, map[string]string{"key": "value"})
	assert.Equal(t, "key: value\n", buf.String())
}

func TestObject(t *testing.T) {
	var buf bytes.Buffer
	print.Object(&buf, "json", map[string]string{"key": "value"})
	exp := `{
	"key": "value"
}
`
	assert.Equal(t, exp, buf.String())

	buf.Reset()
	print.Object(&buf, "yaml", map[string]string{"key": "value"})
	assert.Equal(t, "key: value\n", buf.String())
}

func TestPrint(t *testing.T) {
	var buf bytes.Buffer

	print.Print(&buf, []string{"value1", "value2"})
	assert.Equal(t, "value1\nvalue2\n", buf.String())

	buf.Reset()
	print.Print(&buf, map[string]string{"key": "value"})
	exp := `┌─────┬───────┐
│ KEY │ VALUE │
├─────┼───────┤
│ key │ value │
└─────┴───────┘

`
	assert.Equal(t, exp, buf.String())

	obj := &testStruct{Name: "John", Age: 30}
	buf.Reset()
	print.Print(&buf, obj)
	exp = `{
	"Name": "John",
	"Age": 30
}
`
	assert.Equal(t, exp, buf.String())

	buf.Reset()
	print.Object(&buf, "table", obj)
	assert.Equal(t, exp, buf.String())

	buf.Reset()
	print.Print(&buf, obj)
	assert.Equal(t, exp, buf.String())

	buf.Reset()
	obj2 := struct {
		Name string
		Age  int
	}{Name: "Jane", Age: 25}

	print.Object(&buf, "table", obj2)
	exp = `{
	"Name": "Jane",
	"Age": 25
}
`
	assert.Equal(t, exp, buf.String())

	buf.Reset()
	print.Print(&buf, obj2)
	assert.Equal(t, exp, buf.String())
}

func TestStrings(t *testing.T) {
	var buf bytes.Buffer
	print.Strings(&buf, []string{"value1", "value2"})
	assert.Equal(t, "value1\nvalue2\n", buf.String())
}

func TestMap(t *testing.T) {
	var buf bytes.Buffer
	print.Map(&buf, []string{"Key", "Value"}, map[string]string{"key": "value"})
	exp := `┌─────┬───────┐
│ KEY │ VALUE │
├─────┼───────┤
│ key │ value │
└─────┴───────┘

`
	assert.Equal(t, exp, buf.String())
}

func TestRegisterType(t *testing.T) {
	obj := []*testStruct{
		{Name: "John", Age: 30},
		{Name: "Jane", Age: 25},
	}
	var buf bytes.Buffer
	print.Print(&buf, obj)
	exp := `[
	{
		"Name": "John",
		"Age": 30
	},
	{
		"Name": "Jane",
		"Age": 25
	}
]
`
	assert.Equal(t, exp, buf.String())

	print.RegisterType(([]*testStruct)(nil), func(w io.Writer, value any) {
		print.Yaml(w, value)
	})

	buf.Reset()
	print.Print(&buf, obj)
	exp = `- name: John
  age: 30
- name: Jane
  age: 25
`
	assert.Equal(t, exp, buf.String())

	buf.Reset()
	print.Object(&buf, "table", obj)
	assert.Equal(t, exp, buf.String())
}

type testStruct struct {
	Name string
	Age  int
}

func (t *testStruct) Print(w io.Writer) {
	print.JSON(w, t)
}

func Test_DocumentationOneLine(t *testing.T) {
	w := bytes.NewBuffer([]byte{})

	doc := `line 1
with continuation
	
Line 2.

Line 3
`
	print.TextOneLine(w, doc)
	assert.Equal(t, `line 1 with continuation. Line 2. Line 3.`, w.String())

	exp := `  line 1
  with continuation
  	
  Line 2.
  
  Line 3
  
`
	w.Reset()
	print.Text(w, doc, "  ", false)
	assert.Equal(t, exp, w.String())

	exp = `line 1
  with continuation
  	
  Line 2.
  
  Line 3
  
`
	w.Reset()
	print.Text(w, doc, "  ", true)
	assert.Equal(t, exp, w.String())
}
