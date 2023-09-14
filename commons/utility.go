package commons

import (
	"fmt"
	"strconv"
	"strings"
)

// ConvertInterfaceToInt mengubah tipe data interface{} ke int
func ConvertInterfaceToInt(data interface{}) (int, error) {
	switch value := data.(type) {
	case int:
		return value, nil
	case int64, int32:
		return int(value.(int64)), nil
	case float64, float32:
		return int(value.(float64)), nil
	case string:
		return extractNumberFromString(value)
	default:
		return 0, fmt.Errorf("unsupported type: %T", value)
	}
}

// ConvertInterfaceToInt64 mengubah tipe data interface{} ke int64
func ConvertInterfaceToInt64(data interface{}) (int64, error) {
	switch value := data.(type) {
	case int, int64, int32:
		return int64(value.(int)), nil
	case float64, float32:
		return int64(value.(float64)), nil
	case string:
		intValue, err := extractNumberFromString(value)
		if err != nil {
			return 0, err
		}
		return int64(intValue), nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", value)
	}
}

func extractNumberFromString(str string) (int, error) {
	cleanStr := strings.Trim(str, "{}")
	return strconv.Atoi(cleanStr)
}
