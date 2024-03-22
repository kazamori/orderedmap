package orderedmap_test

import (
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
