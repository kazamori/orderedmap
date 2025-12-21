package orderedmap_test

import (
	"encoding/json"
	"fmt"

	"github.com/kazamori/orderedmap"
)

func Example() {
	m := orderedmap.New[string, any]()
	m.Set("key1", "value1")
	m.Set("key2", 3)
	m.Set("key3", []float64{1.41421356, 3.14})
	for _, v := range m.Pairs() {
		fmt.Println(v.String())
	}
	fmt.Println("================================")
	fmt.Println(m.String())
	// Output:
	// (key1, value1)
	// (key2, 3)
	// (key3, [1.41421356 3.14])
	// ================================
	// {"key1":"value1","key2":3,"key3":[1.41421356,3.14]}
}

func Example_get() {
	m := orderedmap.New[string, int]()
	m.Set("one", 1)
	m.Set("two", 2)
	m.Set("three", 3)

	if val, ok := m.Get("two"); ok {
		fmt.Println("Found:", val)
	}
	if _, ok := m.Get("missing"); !ok {
		fmt.Println("Key not found")
	}
	// Output:
	// Found: 2
	// Key not found
}

func Example_jsonUnmarshal() {
	jsonData := []byte(`{"z":"last","a":"first","m":"middle"}`)
	m := orderedmap.New[string, any]()
	if err := json.Unmarshal(jsonData, m); err != nil {
		panic(err)
	}

	// Keys are preserved in original JSON order
	for _, p := range m.Pairs() {
		fmt.Printf("%s: %v\n", p.Key, p.Value)
	}
	// Output:
	// z: last
	// a: first
	// m: middle
}

func Example_jsonMarshal() {
	m := orderedmap.New[string, any]()
	m.Set("name", "Alice")
	m.Set("age", 30)
	m.Set("active", true)

	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output:
	// {"name":"Alice","age":30,"active":true}
}

func Example_nestedJSON() {
	jsonData := []byte(`{"user":{"name":"Bob","score":100},"items":["a","b","c"]}`)
	m := orderedmap.New[string, any]()
	if err := json.Unmarshal(jsonData, m); err != nil {
		panic(err)
	}

	fmt.Println(m.String())
	// Output:
	// {"user":{"name":"Bob","score":100},"items":["a","b","c"]}
}

func Example_withCapacity() {
	// Pre-allocate capacity for better performance
	m := orderedmap.WithCapacity[string, int](100)
	m.Set("first", 1)
	m.Set("second", 2)
	fmt.Println(m.String())
	// Output:
	// {"first":1,"second":2}
}

func Example_toMap() {
	om := orderedmap.New[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)

	// Convert to standard Go map
	m := orderedmap.ToMap[map[string]int](om)
	fmt.Println(m["a"], m["b"])
	// Output:
	// 1 2
}
