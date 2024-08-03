package objects

import (
	"errors"
	"reflect"
	"testing"
)

type DummyStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type DummyStructA struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
}

func TestMarshalStructToMap(t *testing.T) {
	tests := []struct {
		input    any
		expected map[string]any
		err      error
	}{
		{
			input:    DummyStruct{Name: "Alice", Age: 3},
			expected: map[string]any{"name": "Alice", "age": 3},
			err:      nil,
		},
		{
			input:    struct{}{},
			expected: map[string]any{},
			err:      nil,
		},
		{
			input:    nil,
			expected: nil,
			err:      nil,
		},
	}

	for _, test := range tests {
		result, err := MarshalStructToMap(test.input)
		if (err != nil && test.err == nil) || (err == nil && test.err != nil) || (err != nil && test.err != nil && err.Error() != test.err.Error()) {
			t.Errorf("MarshalStructToMap(%v) error = %v; want %v", test.input, err, test.err)
		}
		if !compareMaps(result, test.expected) {
			t.Errorf("MarshalStructToMap(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestMarshalMapToStruct(t *testing.T) {
	tests := []struct {
		input    map[string]any
		target   any
		expected any
		err      error
	}{
		{
			input:    map[string]any{"name": "Bob", "age": 25},
			target:   &DummyStruct{},
			expected: &DummyStruct{Name: "Bob", Age: 25},
			err:      nil,
		},
		{
			input:    map[string]any{"id": 1, "code": "XYZ"},
			target:   &DummyStructA{},
			expected: &DummyStructA{ID: 1, Code: "XYZ"},
			err:      nil,
		},
		{
			input:    nil,
			target:   &DummyStruct{},
			expected: &DummyStruct{},
			err:      ErrTargetCannotBeNil,
		},
		{
			input:    map[string]any{"id": 1, "code": "XYZ"},
			target:   nil,
			expected: nil,
			err:      ErrDestCannotBeNil,
		},
	}

	for _, test := range tests {
		err := MarshalMapToStruct(test.input, test.target)
		if (err != nil && test.err == nil) || (err == nil && test.err != nil) || (err != nil && test.err != nil && err.Error() != test.err.Error()) {
			t.Errorf("MarshalMapToStruct(%v, %v) error = %v; want %v", test.input, test.target, err, test.err)
		}
		if !compareStructs(test.target, test.expected) {
			t.Errorf("MarshalMapToStruct(%v, %v) = %v; want %v", test.input, test.target, test.target, test.expected)
		}
	}
}

func TestMarshalStructMerge(t *testing.T) {

	type F struct {
		A int `json:"a"`
		B int `json:"b"`
		C int `json:"c"`
		D int `json:"d"`
	}

	type Y struct {
		C int `json:"c"`
		D int `json:"d"`
	}

	type N struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	{
		input := &F{}

		data := &Y{C: 9, D: 1}

		err := MarshalStructMerge(input, data)

		wanted := &F{C: 9, D: 1}

		var expectedErr error = nil

		if err != nil {
			t.Errorf("MarshalStructMerge(%v, %v) error = %v; want %v", input, data, err, expectedErr)
		}

		if !compareStructs(input, wanted) {
			t.Errorf("MarshalStructMerge(%v, %v) = %v; want %v", input, data, input, wanted)
		}

	}

	{
		input := &F{}

		data := []any{
			&Y{C: 2, D: 2},
			&Y{C: 3},
			map[string]any{"a": 1, "b": 18},
		}

		expected := &F{A: 1, B: 18, C: 3, D: 2}

		err := MarshalStructMerge(input, data...)

		var expectedErr error = nil

		if err != nil {
			t.Errorf("MarshalStructMerge(%v, %v) error = %v; want %v", input, data, err, expectedErr)
		}

		if !compareStructs(input, expected) {
			t.Errorf("MarshalStructMerge(%v, %v) = %v; want %v", input, data, input, expected)
		}

	}

	{
		var input any = nil

		data := []any{
			&Y{C: 2, D: 2},
			&N{A: 1, B: 18},
		}

		err := MarshalStructMerge(input, data...)

		var expectedErr error = ErrTargetMustBeAStruct

		if err != nil && !errors.Is(err, expectedErr) {
			t.Errorf("MarshalStructMerge(%v, %v) error = %v; want %v", input, data, err, expectedErr)
		}

	}

}

func compareMaps(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, av := range a {
		bv, ok := b[k]
		if !ok {
			return false
		}
		if !compareValues(av, bv) {
			return false
		}
	}
	return true
}

func compareValues(a, b any) bool {
	switch av := a.(type) {
	case float64:
		if bv, ok := b.(float64); ok {
			return av == bv
		}
		if bi, ok := b.(int); ok {
			return av == float64(bi)
		}
	case int:
		if bv, ok := b.(float64); ok {
			return float64(av) == bv
		}
		if bi, ok := b.(int); ok {
			return av == bi
		}
	default:
		return reflect.DeepEqual(a, b)
	}
	return false
}

// Helper function to compare structs
func compareStructs(a, b any) bool {
	return reflect.DeepEqual(a, b)
}
