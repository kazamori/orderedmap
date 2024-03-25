package orderedmap_test

import (
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
