package values

import (
	"fmt"
)

func ExampleShrink() {
	// Example with nested structures
	input := map[string]any{
		"users": []any{
			map[string]any{"name": "", "age": 0, "active": false},
			map[string]any{"name": "John", "age": 30, "active": true, "tags": []any{"", "admin"}},
			map[string]any{"name": "Jane", "age": 25, "tags": []any{"", "", ""}},
		},
		"config": map[string]any{
			"debug":   false,
			"timeout": 0,
			"host":    "localhost",
		},
		"empty_section": map[string]any{
			"a": "",
			"b": 0,
			"c": []any{},
		},
	}

	result := Shrink(input)
	fmt.Printf("Result: %+v\n", result)

	// Example with simple values
	fmt.Printf("Empty string: %v\n", Shrink(""))
	fmt.Printf("Non-empty string: %v\n", Shrink("hello"))
	fmt.Printf("Zero int: %v\n", Shrink(0))
	fmt.Printf("Non-zero int: %v\n", Shrink(42))
	fmt.Printf("Empty slice: %v\n", Shrink([]any{}))
	fmt.Printf("Slice with empty values: %v\n", Shrink([]any{"", 0, false}))
	fmt.Printf("Slice with mixed values: %v\n", Shrink([]any{"", "hello", 0, 42}))

	// Output:
	// Result: map[config:map[host:localhost] users:[map[active:true age:30 name:John tags:[admin]] map[age:25 name:Jane]]]
	// Empty string: <nil>
	// Non-empty string: hello
	// Zero int: <nil>
	// Non-zero int: 42
	// Empty slice: <nil>
	// Slice with empty values: <nil>
	// Slice with mixed values: [hello 42]
}

func ExampleMapAny_Shrink() {
	// Example with MapAny
	m := MapAny{
		"empty":     "",
		"text":      "hello",
		"zero":      0,
		"num":       42,
		"false":     false,
		"true":      true,
		"nil":       nil,
		"empty_map": MapAny{"a": "", "b": 0},
		"valid_map": MapAny{"x": "world", "y": 100},
	}

	result := m.Shrink()
	fmt.Printf("Result: %+v\n", result)

	// Example with all empty values
	empty := MapAny{"a": "", "b": 0, "c": false, "d": nil}
	emptyResult := empty.Shrink()
	if emptyResult == nil {
		fmt.Printf("Empty result: <nil>\n")
	} else {
		fmt.Printf("Empty result: %v\n", emptyResult)
	}

	// Output:
	// Result: map[num:42 text:hello true:true valid_map:map[x:world y:100]]
	// Empty result: <nil>
}
