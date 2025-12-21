# orderedmap

A generic ordered map implementation in Go that preserves insertion order and supports JSON serialization/deserialization.

## Installation

```bash
go get github.com/kazamori/orderedmap
```

## Features

- Generic type support with `OrderedMap[K, V]`
- Preserves insertion order of keys
- JSON marshaling/unmarshaling with order preservation
- O(1) key lookup

## Quick Start

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/kazamori/orderedmap"
)

func main() {
    // Create a new ordered map
    m := orderedmap.New[string, any]()
    m.Set("name", "Alice")
    m.Set("age", 30)
    m.Set("active", true)

    // Iterate in insertion order
    for _, p := range m.Pairs() {
        fmt.Printf("%s: %v\n", p.Key, p.Value)
    }

    // Get a value
    if val, ok := m.Get("name"); ok {
        fmt.Println("Name:", val)
    }

    // Serialize to JSON (preserves order)
    data, _ := json.Marshal(m)
    fmt.Println(string(data))
    // Output: {"name":"Alice","age":30,"active":true}
}
```

## JSON Unmarshaling

Parse JSON while preserving the original key order:

```go
jsonData := []byte(`{"z":"last","a":"first","m":"middle"}`)
m := orderedmap.New[string, any]()
json.Unmarshal(jsonData, m)

for _, p := range m.Pairs() {
    fmt.Printf("%s: %v\n", p.Key, p.Value)
}
// Output:
// z: last
// a: first
// m: middle
```

## API Reference

| Function/Method | Description |
|----------------|-------------|
| `New[K, V]()` | Create a new empty ordered map |
| `WithCapacity[K, V](size)` | Create with pre-allocated capacity |
| `NewFromMap[M, K, V](m)` | Create from a standard Go map |
| `Set(key, value)` | Add or update a key-value pair |
| `Get(key)` | Retrieve a value by key |
| `Pairs()` | Get all pairs in insertion order |
| `String()` | JSON string representation |
| `ToMap[M, K, V](om)` | Convert to standard Go map |

## Examples

See [example_test.go](./example_test.go) for more usage examples.

## CLI Tool

A CLI tool is included for testing JSON round-trip serialization:

```bash
make build
./bin/cli -data '{"s":"test","i":3,"a":[{"f":3.14},{"b":true}]}'
```

Output:
```json
{
  "s": "test",
  "i": 3,
  "a": [
    {
      "f": 3.14
    },
    {
      "b": true
    }
  ]
}
```
