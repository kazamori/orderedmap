package orderedmap

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
)

func isSlice(v any) bool {
	if !reflect.ValueOf(v).IsValid() {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

func isValueAndTypeAreSlice(typ any, value any) bool {
	return isSlice(typ) && isSlice(value)
}

func convertValuesForSlice(typ any, anyValues []any) any {
	elem := reflect.TypeOf(typ).Elem().Elem()
	switch kind := elem.Kind(); kind {
	case reflect.Bool:
		values := make([]bool, 0, len(anyValues))
		for _, v := range anyValues {
			values = append(values, v.(bool))
		}
		return values
	case reflect.Int:
		values := make([]int, 0, len(anyValues))
		for _, v := range anyValues {
			values = append(values, v.(int))
		}
		return values
	case reflect.Int64:
		values := make([]int64, 0, len(anyValues))
		for _, v := range anyValues {
			values = append(values, v.(int64))
		}
		return values
	case reflect.Float64:
		values := make([]float64, 0, len(anyValues))
		for _, v := range anyValues {
			values = append(values, v.(float64))
		}
		return values
	case reflect.String:
		values := make([]string, 0, len(anyValues))
		for _, v := range anyValues {
			values = append(values, v.(string))
		}
		return values

	default:
		panic(fmt.Errorf("unsupported type to convert: %s", kind))
	}
}

func isObject(delim json.Delim) bool {
	return delim == leftCurlyBrace
}

func decodeKey(decoder *json.Decoder) (string, error) {
	token, err := decoder.Token()
	if err != nil {
		return "", fmt.Errorf("failed to get token as key: %w", err)
	}
	if delim, ok := token.(json.Delim); ok {
		switch delim {
		case leftCurlyBrace:
			return "", ErrNestedObject
		case leftBracket:
			return "", ErrNestedArray
		case rightCurlyBrace, rightBracket:
			return "", ErrEndOfJSON
		}
	}
	if key, ok := token.(string); ok {
		return key, nil
	}
	return "", fmt.Errorf("key should be string, but not %v", token)
}

func decodeObject(decoder *json.Decoder) (any, error) {
	m := New[string, any]()
	for {
		key, err := decodeKey(decoder)
		if err != nil {
			if errors.Is(err, ErrEndOfJSON) {
				return m, nil
			} else if errors.Is(err, ErrNestedObject) {
				return decodeObject(decoder)
			} else if errors.Is(err, ErrNestedArray) {
				return decodeArray(decoder)
			}
			return nil, fmt.Errorf("failed to get key: %w", err)
		}
		value, err := decodeValue(decoder)
		if err != nil {
			return nil, fmt.Errorf("failed to get value: %w", err)
		}
		m.Set(key, value)
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
	switch t := token.(type) {
	case string, float64, int64, bool:
		return token, nil
	default:
		slog.Debug("unexpected token type", "token", t)
	}
	return token, nil
}

func decodeValue(decoder *json.Decoder) (any, error) {
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
