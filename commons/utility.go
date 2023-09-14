package commons

import (
	"errors"
	"fmt"
	"strconv"
)

func InterfaceToInt(data interface{}) (int, error) {
	switch value := data.(type) {
	case int:
		return value, nil
	case int32:
		return int(value), nil
	case int64:
		return int(value), nil
	case float32:
		return int(value), nil
	case float64:
		return int(value), nil
	case string:
		id, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("failed converting string to int: %w", err)
		}
		return id, nil
	default:
		return 0, errors.New("unsupported type for conversion to int")
	}
}

func InterfaceToInt64(data interface{}) (int64, error) {
	intValue, err := InterfaceToInt(data)
	if err != nil {
		return 0, err
	}
	return int64(intValue), nil
}
