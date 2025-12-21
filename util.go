package orderedmap

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

func isSlice(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Slice
}

func assertType[T any](value any) (T, error) {
	var zero T
	if isSlice(zero) {
		v, ok := value.([]any)
		if !ok {
			return zero, fmt.Errorf("failed to assert %v (%T) as slice", value, value)
		}
		values, err := convertValuesForSlice(zero, v)
		if err != nil {
			return zero, fmt.Errorf("failed to convert slice: %w", err)
		}
		typedValues, ok := any(values).(T)
		if !ok {
			return zero, fmt.Errorf("failed to convert typed slice (%T)", values)
		}
		return typedValues, nil
	}
	v, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("failed to assert %v (%T) as %T", value, value, zero)
	}
	return v, nil
}

func convertValuesForSlice(zero any, anyValues []any) (any, error) {
	if zero == nil {
		return nil, errors.New("zero cannot be nil")
	}

	elem := reflect.TypeOf(zero).Elem()
	switch kind := elem.Kind(); kind {
	case reflect.Bool:
		values := make([]bool, 0, len(anyValues))
		for _, v := range anyValues {
			vv, ok := v.(bool)
			if !ok {
				return nil, fmt.Errorf("failed to assert %v (%T) as bool", v, v)
			}
			values = append(values, vv)
		}
		return values, nil
	case reflect.Int:
		values := make([]int, 0, len(anyValues))
		for _, v := range anyValues {
			vv, ok := v.(int)
			if !ok {
				return nil, fmt.Errorf("failed to assert %v (%T) as int", v, v)
			}
			values = append(values, vv)
		}
		return values, nil
	case reflect.Int64:
		values := make([]int64, 0, len(anyValues))
		for _, v := range anyValues {
			vv, ok := v.(int64)
			if !ok {
				return nil, fmt.Errorf("failed to assert %v (%T) as int64", v, v)
			}
			values = append(values, vv)
		}
		return values, nil
	case reflect.Uint:
		values := make([]uint, 0, len(anyValues))
		for _, v := range anyValues {
			vv, ok := v.(uint)
			if !ok {
				return nil, fmt.Errorf("failed to assert %v (%T) as uint", v, v)
			}
			values = append(values, vv)
		}
		return values, nil
	case reflect.Uint64:
		values := make([]uint64, 0, len(anyValues))
		for _, v := range anyValues {
			vv, ok := v.(uint64)
			if !ok {
				return nil, fmt.Errorf("failed to assert %v (%T) as uint64", v, v)
			}
			values = append(values, vv)
		}
		return values, nil
	case reflect.Float64:
		values := make([]float64, 0, len(anyValues))
		for _, v := range anyValues {
			vv, ok := v.(float64)
			if !ok {
				return nil, fmt.Errorf("failed to assert %v (%T) as float64", v, v)
			}
			values = append(values, vv)
		}
		return values, nil
	case reflect.String:
		values := make([]string, 0, len(anyValues))
		for _, v := range anyValues {
			vv, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("failed to assert %v (%T) as string", v, v)
			}
			values = append(values, vv)
		}
		return values, nil
	default:
		return nil, fmt.Errorf("unsupported type to convert: %s", kind)
	}
}

func isObject(delim json.Delim) bool {
	return delim == leftCurlyBrace
}

func decodeObject(decoder *json.Decoder) (any, error) {
	m := New[string, any]()
	for {
		if err := m.decodeKeyAndValue(decoder); err != nil {
			if errors.Is(err, ErrEndOfJSON) {
				return m, nil
			}
			return nil, fmt.Errorf("failed to decode key/value: %w", err)
		}
	}
}

func decodeArray(decoder *json.Decoder) ([]any, error) {
	values := make([]any, 0)
	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("failed to get token for an element: %w", err)
		}
		if delim, ok := token.(json.Delim); ok {
			switch delim {
			case leftCurlyBrace:
				nestedObject, err := decodeObject(decoder)
				if err != nil {
					return nil, fmt.Errorf("failed to decode nested object: %w", err)
				}
				values = append(values, nestedObject)
			case leftBracket:
				nestedValues, err := decodeArray(decoder)
				if err != nil {
					return nil, fmt.Errorf("failed to decode nested array: %w", err)
				}
				values = append(values, nestedValues)
			case rightBracket:
				return values, nil
			default:
				return nil, fmt.Errorf("unsupported format: %s", delim)
			}
		} else {
			value, err := handleToken(token)
			if err != nil {
				return nil, fmt.Errorf("failed to decode value for an element: %w", err)
			}
			values = append(values, value)
		}
	}
}

func handleToken(token json.Token) (any, error) {
	// TODO: An abstraction layer for handling a token if it needed
	if token == nil {
		return token, nil
	}
	switch t := token.(type) {
	case string, float64, int, int64, uint, uint64, bool:
		return token, nil
	default:
		return nil, fmt.Errorf("unexpected token type: %v", t)
	}
}

func decodeValueAsAny(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to get token as value: %w", err)
	}
	if delim, ok := token.(json.Delim); ok {
		switch delim {
		case leftCurlyBrace:
			return decodeObject(decoder)
		case leftBracket:
			return decodeArray(decoder)
		default:
			return nil, fmt.Errorf("unsupported format: %s", delim)
		}
	}
	return handleToken(token)
}
