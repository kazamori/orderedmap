package orderedmap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	colon = []byte{':'}
	comma = []byte{','}
)

var (
	leftCurlyBrace  = json.Delim('{')
	rightCurlyBrace = json.Delim('}')
	leftBracket     = json.Delim('[')
	rightBracket    = json.Delim(']')
)

var (
	ErrEndOfJSON    = errors.New("end of JSON")
	ErrNestedObject = errors.New("detect nested object")
	ErrNestedArray  = errors.New("detect nested array")
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
	buf.Write([]byte(leftCurlyBrace.String()))
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
	buf.Write([]byte(rightCurlyBrace.String()))
	return buf.Bytes(), nil
}

func (m *OrderedMap[K, V]) initialize() {
	m.pairs = make([]Pair[K, V], 0)
	m.index = make(map[K]int)
	m.pos = 0
}

func (m *OrderedMap[K, V]) decodeKeyAndValue(decoder *json.Decoder) error {
	key, err := decodeKey(decoder)
	if err != nil {
		return fmt.Errorf("failed to get key: %w", err)
	}
	value, err := decodeValue(decoder)
	if err != nil {
		return fmt.Errorf("failed to get value: %w", err)
	}
	var v V
	if isValueAndTypeAreSlice(v, value) {
		v = any(convertValuesForSlice(new(V), value.([]any))).(V)
	} else {
		v = value.(V)
	}
	m.Set(any(key).(K), v)
	return nil
}

func (m *OrderedMap[K, V]) UnmarshalJSON(b []byte) error {
	m.initialize()
	decoder := json.NewDecoder(bytes.NewReader(b))
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	if delim, ok := token.(json.Delim); !(ok && isObject(delim)) {
		return fmt.Errorf("expects JSON object, but not: %s", delim)
	}
	for {
		if err := m.decodeKeyAndValue(decoder); err != nil {
			if errors.Is(err, ErrEndOfJSON) {
				return nil
			}
			return fmt.Errorf("failed to decode key/value: %w", err)
		}
	}
}

func ToMap[M ~map[K]V, K comparable, V any](om *OrderedMap[K, V]) M {
	pairs := om.Pairs()
	m := make(M, len(pairs))
	for _, p := range pairs {
		m[p.Key] = p.Value
	}
	return m
}
