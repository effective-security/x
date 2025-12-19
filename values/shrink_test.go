package values

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShrink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "non-empty string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "zero int",
			input:    0,
			expected: nil,
		},
		{
			name:     "non-zero int",
			input:    42,
			expected: 42,
		},
		{
			name:     "false bool",
			input:    false,
			expected: nil,
		},
		{
			name:     "true bool",
			input:    true,
			expected: true,
		},
		{
			name:     "zero float",
			input:    0.0,
			expected: nil,
		},
		{
			name:     "non-zero float",
			input:    3.14,
			expected: 3.14,
		},
		{
			name:     "empty slice",
			input:    []any{},
			expected: nil,
		},
		{
			name:     "slice with only empty values",
			input:    []any{"", 0, false, nil},
			expected: nil,
		},
		{
			name:     "slice with mixed values",
			input:    []any{"", "hello", 0, 42, false, true, nil},
			expected: []any{"hello", 42, true},
		},
		{
			name:     "nested slice",
			input:    []any{[]any{"", 0}, []any{"hello", 42}, []any{}},
			expected: []any{[]any{"hello", 42}},
		},
		{
			name:     "empty map",
			input:    map[string]any{},
			expected: nil,
		},
		{
			name:     "map with only empty values",
			input:    map[string]any{"a": "", "b": 0, "c": false, "d": nil},
			expected: nil,
		},
		{
			name:     "map with mixed values",
			input:    map[string]any{"empty": "", "text": "hello", "zero": 0, "num": 42, "false": false, "true": true, "nil": nil},
			expected: map[string]any{"text": "hello", "num": 42, "true": true},
		},
		{
			name: "nested map",
			input: map[string]any{
				"empty_map": map[string]any{"a": "", "b": 0},
				"valid_map": map[string]any{"x": "hello", "y": 42},
				"mixed_map": map[string]any{"empty": "", "valid": "world"},
			},
			expected: map[string]any{
				"valid_map": map[string]any{"x": "hello", "y": 42},
				"mixed_map": map[string]any{"valid": "world"},
			},
		},
		{
			name: "complex nested structure",
			input: map[string]any{
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
			},
			expected: map[string]any{
				"users": []any{
					map[string]any{"name": "John", "age": 30, "active": true, "tags": []any{"admin"}},
					map[string]any{"name": "Jane", "age": 25},
				},
				"config": map[string]any{
					"host": "localhost",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Shrink(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapAnyShrink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    MapAny
		expected MapAny
	}{
		{
			name:     "nil MapAny",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty MapAny",
			input:    MapAny{},
			expected: nil,
		},
		{
			name:     "MapAny with only empty values",
			input:    MapAny{"a": "", "b": 0, "c": false, "d": nil},
			expected: nil,
		},
		{
			name:     "MapAny with mixed values",
			input:    MapAny{"empty": "", "text": "hello", "zero": 0, "num": 42, "false": false, "true": true, "nil": nil},
			expected: MapAny{"text": "hello", "num": 42, "true": true},
		},
		{
			name: "MapAny with nested structures",
			input: MapAny{
				"users": []any{
					map[string]any{"name": "", "active": false},
					map[string]any{"name": "John", "active": true},
				},
				"empty_list": []any{"", 0, false},
				"valid_data": "test",
			},
			expected: MapAny{
				"users": []any{
					map[string]any{"name": "John", "active": true},
				},
				"valid_data": "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.input.Shrink()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{
			name:     "nil",
			input:    nil,
			expected: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "non-empty string",
			input:    "hello",
			expected: false,
		},
		{
			name:     "zero int",
			input:    0,
			expected: true,
		},
		{
			name:     "non-zero int",
			input:    42,
			expected: false,
		},
		{
			name:     "zero int64",
			input:    int64(0),
			expected: true,
		},
		{
			name:     "non-zero int64",
			input:    int64(42),
			expected: false,
		},
		{
			name:     "zero uint",
			input:    uint(0),
			expected: true,
		},
		{
			name:     "non-zero uint",
			input:    uint(42),
			expected: false,
		},
		{
			name:     "zero float32",
			input:    float32(0),
			expected: true,
		},
		{
			name:     "non-zero float32",
			input:    float32(3.14),
			expected: false,
		},
		{
			name:     "zero float64",
			input:    0.0,
			expected: true,
		},
		{
			name:     "non-zero float64",
			input:    3.14,
			expected: false,
		},
		{
			name:     "false bool",
			input:    false,
			expected: true,
		},
		{
			name:     "true bool",
			input:    true,
			expected: false,
		},
		{
			name:     "empty slice",
			input:    []any{},
			expected: true,
		},
		{
			name:     "non-empty slice",
			input:    []any{1, 2, 3},
			expected: false,
		},
		{
			name:     "empty map",
			input:    map[string]any{},
			expected: true,
		},
		{
			name:     "non-empty map",
			input:    map[string]any{"key": "value"},
			expected: false,
		},
		{
			name:     "nil pointer",
			input:    (*string)(nil),
			expected: true,
		},
		{
			name:     "non-nil pointer to empty string",
			input:    func() *string { s := ""; return &s }(),
			expected: true,
		},
		{
			name:     "non-nil pointer to non-empty string",
			input:    func() *string { s := "hello"; return &s }(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := IsEmpty(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestShrinkEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("deeply nested empty structures", func(t *testing.T) {
		t.Parallel()
		input := map[string]any{
			"level1": map[string]any{
				"level2": map[string]any{
					"level3": []any{
						map[string]any{
							"empty": "",
							"zero":  0,
						},
					},
				},
			},
		}
		result := Shrink(input)
		assert.Nil(t, result)
	})

	t.Run("mixed types in slice", func(t *testing.T) {
		t.Parallel()
		input := []any{
			"",
			0,
			false,
			[]any{},
			map[string]any{},
			"valid",
			42,
			true,
			[]any{"nested"},
			map[string]any{"key": "value"},
		}
		expected := []any{
			"valid",
			42,
			true,
			[]any{"nested"},
			map[string]any{"key": "value"},
		}
		result := Shrink(input)
		assert.Equal(t, expected, result)
	})

	t.Run("slice with all empty nested structures", func(t *testing.T) {
		t.Parallel()
		input := []any{
			[]any{},
			map[string]any{},
			[]any{"", 0, false},
			map[string]any{"a": "", "b": 0},
		}
		result := Shrink(input)
		assert.Nil(t, result)
	})
}

func TestShrinkPreservesTypes(t *testing.T) {
	t.Parallel()

	t.Run("preserves non-empty primitive types", func(t *testing.T) {
		t.Parallel()

		// Test different numeric types
		assert.Equal(t, int8(42), Shrink(int8(42)))
		assert.Equal(t, int16(42), Shrink(int16(42)))
		assert.Equal(t, int32(42), Shrink(int32(42)))
		assert.Equal(t, int64(42), Shrink(int64(42)))
		assert.Equal(t, uint8(42), Shrink(uint8(42)))
		assert.Equal(t, uint16(42), Shrink(uint16(42)))
		assert.Equal(t, uint32(42), Shrink(uint32(42)))
		assert.Equal(t, uint64(42), Shrink(uint64(42)))
		assert.Equal(t, float32(3.14), Shrink(float32(3.14)))
		assert.Equal(t, float64(3.14), Shrink(float64(3.14)))
	})

	t.Run("returns nil for empty primitive types", func(t *testing.T) {
		t.Parallel()

		// Test different numeric types
		assert.Nil(t, Shrink(int8(0)))
		assert.Nil(t, Shrink(int16(0)))
		assert.Nil(t, Shrink(int32(0)))
		assert.Nil(t, Shrink(int64(0)))
		assert.Nil(t, Shrink(uint8(0)))
		assert.Nil(t, Shrink(uint16(0)))
		assert.Nil(t, Shrink(uint32(0)))
		assert.Nil(t, Shrink(uint64(0)))
		assert.Nil(t, Shrink(float32(0)))
		assert.Nil(t, Shrink(float64(0)))
	})
}
