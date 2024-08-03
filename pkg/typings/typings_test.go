package typings

import (
	"testing"
)

func TestIsMap(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: false,
		},
		{
			name:     "valid map",
			input:    map[string]int{"one": 1, "two": 2},
			expected: true,
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: true,
		},
		{
			name:     "non-map input (slice)",
			input:    []int{1, 2, 3},
			expected: false,
		},
		{
			name:     "non-map input (string)",
			input:    "test string",
			expected: false,
		},
		{
			name:     "non-map input (int)",
			input:    123,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsMap(tt.input)
			if got != tt.expected {
				t.Errorf("IsMap() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

type DummyStruct struct {
	Name string
	Age  int
}

type DummyStructA struct {
	ID   int
	Code string
}

func TestIsPointerToStruct(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		// Pointer to struct
		{&DummyStruct{}, true},
		{&DummyStructA{}, true},

		// Not a pointer to struct
		{DummyStruct{}, false},
		{DummyStructA{}, false},
		{nil, false},
		{42, false},
		{"string", false},
		{map[string]int{}, false},
		{func() {}, false},

		// Pointer to non-struct
		{(*int)(nil), false},
		{(*string)(nil), false},
	}

	for _, test := range tests {
		result := IsPointerToStruct(test.value)
		if result != test.expected {
			t.Errorf("IsPointerToStruct(%v) = %v; want %v", test.value, result, test.expected)
		}
	}
}

// Test the IsStruct function
func TestIsStruct(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		// Struct
		{DummyStruct{}, true},
		{DummyStructA{}, true},

		// Not a struct
		{nil, false},
		{42, false},
		{"string", false},
		{map[string]int{}, false},
		{func() {}, false},
		{&DummyStruct{}, false},  // Pointer to struct is not a struct itself
		{&DummyStructA{}, false}, // Pointer to struct is not a struct itself
	}

	for _, test := range tests {
		result := IsStruct(test.value)
		if result != test.expected {
			t.Errorf("IsStruct(%v) = %v; want %v", test.value, result, test.expected)
		}
	}
}

func TestIsZeroValue(t *testing.T) {

	tests := []struct {
		value    interface{}
		expected bool
	}{
		// Nil cases
		{nil, true},
		{(*DummyStruct)(nil), true},
		{map[string]int(nil), true},
		{chan int(nil), true},
		{interface{}(nil), true},

		// Zero cases
		{"", true},
		{0, true},
		{0.0, true},
		{complex(0, 0), true},
		{false, true},

		// Non-zero cases
		{"Hello", false},
		{42, false},
		{3.14, false},
		{complex(1, 1), false},
		{true, false},

		// Struct cases
		{DummyStruct{Name: "John", Age: 30}, false},
		{&DummyStruct{Name: "John", Age: 30}, false},
	}

	for _, test := range tests {
		result := IsZeroValue(test.value)
		if result != test.expected {
			t.Errorf("isNotNilOrZero(%v) = %v; want %v", test.value, result, test.expected)
		}
	}
}
