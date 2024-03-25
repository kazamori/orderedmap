package orderedmap_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kazamori/orderedmap"
)

func TestOrderedMap(t *testing.T) {
	tests := []struct {
		name     string
		pairs    []orderedmap.Pair[string, any]
		expected string
	}{
		{
			name: "simple",
			pairs: []orderedmap.Pair[string, any]{
				{"s", "test"},
				{"i2", 5},
				{"s", "test2"},
				{"f", 3.14},
				{"i", 3},
				{"s", "test3"},
				{"b2", false},
				{"i3", 1},
				{"i", 13},
				{"b1", true},
			},
			expected: `{"s":"test3","i2":5,"f":3.14,"i":13,"b2":false,"i3":1,"b1":true}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := orderedmap.New[string, any]()
			for _, p := range tt.pairs {
				m.Set(p.Key, p.Value)
			}
			actual := m.String()
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Error(diff)
				return
			}
		})
	}
}

func TestFromMap(t *testing.T) {
	tests := []struct {
		name     string
		m        map[string]any
		expected string
	}{
		{
			name: "an int value in a map",
			m: map[string]any{
				"key": 1,
			},
			expected: `{"key":1}`,
		},
		{
			name: "complex data in an array in a map",
			m: map[string]any{
				"a": []map[string]any{
					map[string]any{
						"k1": "test",
					},
					map[string]any{
						"k2": 3.14,
					},
					map[string]any{
						"k3": map[string]any{
							"k3-1": 33,
						},
					},
				},
			},
			expected: `{"a":[{"k1":"test"},{"k2":3.14},{"k3":{"k3-1":33}}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			om := orderedmap.NewFromMap(tt.m)
			actual := om.String()
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Error(diff)
				return
			}
		})
	}
}

func TestToMap(t *testing.T) {
	nested1 := orderedmap.New[string, any]()
	nested1.Set("f", 3.14)
	nested1.Set("s", "nested")
	nested1.Set("b", true)

	nested2 := orderedmap.New[string, any]()
	nested2.Set("f2", 3.14)
	nested2.Set("s2", "nested")
	nested2.Set("b2", true)
	nested2.Set("a2", []string{"v1", "v2", "v3"})

	tests := []struct {
		name     string
		om       *orderedmap.OrderedMap[string, any]
		expected map[string]any
	}{
		{
			name: "simple",
			om: func() *orderedmap.OrderedMap[string, any] {

				om := orderedmap.New[string, any]()
				om.Set("key1", 1)
				om.Set("key2", 2)
				om.Set("key3", 3)
				om.Set("a1", []string{"v1", "v2", "v3"})
				om.Set("a2", []orderedmap.OrderedMap[string, any]{
					*nested2,
					*nested1,
				})
				return om
			}(),
			expected: map[string]any{
				"a1": []string{"v1", "v2", "v3"},
				"a2": []orderedmap.OrderedMap[string, any]{
					*nested2,
					*nested1,
				},
				"key3": 3,
				"key2": 2,
				"key1": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := orderedmap.ToMap[map[string]any](tt.om)
			opts := []cmp.Option{
				cmp.AllowUnexported(orderedmap.OrderedMap[string, any]{}),
			}
			if diff := cmp.Diff(tt.expected, actual, opts...); diff != "" {
				t.Error(diff)
				return
			}
		})
	}
}

func TestUnmarhalAndMarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{
			name: "primitive only",
			json: `{"s":"test","i":3,"n":9007199254740992,"f":3.14,"b":true}`,
		},
		{
			name: "simple array",
			json: `{"a":["test",3,9007199254740992,3.14,true]}`,
		},
		{
			name: "multiple nested array",
			json: `{"a":[["test",3],[9007199254740992,3.14,true,["nested",5]]]}`,
		},
		{
			name: "simple object",
			json: `{"m":{"s":"test","i":7}}`,
		},
		{
			name: "nested object",
			json: `{"m":{"s":"test","i":7,"m2":{"f":3.14,"b":true,"m":{"i":19}}}}`,
		},
		{
			name: "maps in array",
			json: `{"a":[{"s":"test","i":7},{"s2":"test2","f":3.14},{"b":true}]}`,
		},
		{
			name: "complex json data including arrays and objects",
			json: `{"m":{"s":"test","a":[["test",3],[9007199254740992,3.14,true,["nested",5]]],"i":7,"m2":{"f":3.14,"a":[{"s":"test","i":7},{"s2":"test2","f":3.14},{"b":true}],"b":true,"m":{"a":["test",3,9007199254740992,3.14,true],"i":19}}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var m orderedmap.OrderedMap[string, any]
			if err := json.Unmarshal(json.RawMessage(tt.json), &m); err != nil {
				t.Error(err)
				return
			}
			actual := m.String()
			if diff := cmp.Diff(tt.json, actual); diff != "" {
				t.Error(diff)
				return
			}
		})
	}
}

func TestErrorUnmarhalAndMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected error
	}{
		{
			name:     "simple array",
			json:     `["test",3,9007199254740992,3.14,true]`,
			expected: errors.New("expects JSON object, but not: ["),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var m orderedmap.OrderedMap[string, any]
			err := json.Unmarshal(json.RawMessage(tt.json), &m)
			if err == nil {
				t.Error("expects an error occurred, but no error")
			}
			if diff := cmp.Diff(tt.expected.Error(), err.Error()); diff != "" {
				t.Error(diff)
				return
			}
		})
	}
}
