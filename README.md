# orderedmap
![Coverage](https://img.shields.io/badge/Coverage-65.2%25-yellow)

An alternative generic ordered map in Go with de/serializing from/to JSON.

## How to use

See the [example code](./example_test.go).

## test CLI

For testing, deserialize JSON and then serialize it while preserving the order.

```
$ ./bin/cli -data '{"s":"test","i":3,"a":[{"f":3.14},{"b":true}]}'
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
