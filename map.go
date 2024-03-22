package orderedmap

import (
	"bytes"
	"encoding/json"
	"fmt"
)

var (
	leftCurlyBrace  = []byte{'{'} //lint:ignore U1000 May use in the future
	rightCurlyBrace = []byte{'}'} //lint:ignore U1000 May use in the future
	leftBracket     = []byte{'['} //lint:ignore U1000 May use in the future
	rightBracket    = []byte{']'} //lint:ignore U1000 May use in the future
	colon           = []byte{':'} //lint:ignore U1000 May use in the future
	comma           = []byte{','} //lint:ignore U1000 May use in the future
	dot             = []byte{'.'} //lint:ignore U1000 May use in the future
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func (p *Pair[K, V]) String() string {
	return fmt.Sprintf("(%v, %v)", p.Key, p.Value)
}

func (p *Pair[K, V]) QuotedKey() string {
	// FIXME: how to convert comparable to string
	return fmt.Sprintf("%q", fmt.Sprintf("%v", p.Key))
}

func NewPair[K comparable, V any](k K, v V) *Pair[K, V] {
	return &Pair[K, V]{
		Key:   k,
		Value: v,
	}
}

type OrderedMap[K comparable, V any] struct {
	pairs []Pair[K, V]
	index map[K]int
	pos   int
}

func (m *OrderedMap[K, V]) Set(key K, value V) {
	pair := NewPair(key, value)
	if i, ok := m.index[key]; ok {
		m.pairs[i] = *pair
		return
	}
	m.pairs = append(m.pairs, *NewPair(key, value))
	m.index[key] = m.pos
	m.pos += 1
}

func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	i, ok := m.index[key]
	if !ok {
		return *new(V), ok
	}
	pair := m.pairs[i]
	return pair.Value, ok
}

func (m *OrderedMap[K, V]) Pairs() []Pair[K, V] {
	return m.pairs
}

func (m *OrderedMap[K, V]) String() string {
	b, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("invalid data: %w", err))
	}
	return string(b)
}

func WithCapacity[K comparable, V any](size int) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		pairs: make([]Pair[K, V], 0, size),
		index: make(map[K]int, size),
	}
}

func New[K comparable, V any]() *OrderedMap[K, V] {
	return WithCapacity[K, V](0)
}

func NewFromMap[M ~map[K]V, K comparable, V any](m M) *OrderedMap[K, V] {
	om := WithCapacity[K, V](len(m))
	for k, v := range m {
		om.Set(k, v)
	}
	return om
}

func (m OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Write(leftCurlyBrace)
	pairs := m.Pairs()
	length := len(pairs) - 1
	for i, pair := range pairs {
		v, err := json.Marshal(pair.Value)
		if err != nil {
			return nil, err
		}
		buf.Write([]byte(pair.QuotedKey()))
		buf.Write(colon)
		buf.Write(v)
		if i < length {
			buf.Write(comma)
		}
	}
	buf.Write(rightCurlyBrace)
	return buf.Bytes(), nil
}

func ToMap[M ~map[K]V, K comparable, V any](om *OrderedMap[K, V]) M {
	pairs := om.Pairs()
	m := make(M, len(pairs))
	for _, p := range pairs {
		m[p.Key] = p.Value
	}
	return m
}
