# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library (`github.com/kazamori/orderedmap`) that provides a generic ordered map implementation with JSON serialization/deserialization support. Unlike Go's built-in map, this preserves insertion order of keys.

## Common Commands

```bash
# Run all tests with race detection and coverage
make test

# Run linting (go vet + staticcheck)
make lint

# Build the CLI tool
make build

# Run a single test
go test -run TestOrderedMap ./...

# Run tests with verbose output
go test -v -race -cover ./...

# Clean build artifacts and test cache
make clean
```

## Architecture

The library consists of two main files:

- **map.go**: Core `OrderedMap[K, V]` generic struct with:
  - `pairs []Pair[K, V]` - maintains insertion order
  - `index map[K]int` - provides O(1) key lookup
  - Implements `json.Marshaler` and `json.Unmarshaler` interfaces

- **util.go**: JSON decoding helpers that recursively handle nested objects/arrays while preserving order. Uses `json.Decoder` token-by-token parsing rather than standard unmarshaling.

Key design decisions:
- Nested JSON objects are decoded as `*OrderedMap[string, any]`
- JSON arrays are decoded as `[]any`
- Type assertions with generics handled via reflection in `assertType[T]` and `convertValuesForSlice`

## Dependencies

- Go 1.23+
- `github.com/google/go-cmp` (testing only)
- `staticcheck` for linting (install via `go install honnef.co/go/tools/cmd/staticcheck@latest`)
