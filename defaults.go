package dflag

import (
	"fmt"
	"reflect"
	"strconv"
)

type cliKey string

func (k cliKey) String() string {
	return string(k)
}

var (
	cli   cliKey = "flag"
	usage cliKey = "usage" //optional
	value cliKey = "default"
)

func getDefaultUsage(field reflect.StructField, kind reflect.Kind) string {
	if key := field.Tag.Get(usage.String()); len(key) != 0 {
		return key
	}
	return fmt.Sprintf("enter value for type %s", kind.String())
}

func getDefaultValue(field reflect.StructField, kind reflect.Kind) any {
	if key := field.Tag.Get(value.String()); len(key) != 0 {
		return defaultValueForKind(key, kind)
	}

	return getDefaultValueForKind(kind)
}

func getDefaultValueForString(field reflect.StructField) string {
	return getDefaultValue(field, reflect.String).(string)
}

func getDefaultValueForBool(field reflect.StructField) bool {
	return getDefaultValue(field, reflect.Bool).(bool)
}

func getDefaultValueForInt(field reflect.StructField) int64 {
	return getDefaultValue(field, reflect.Int64).(int64)
}

func getDefaultValueForFloat(field reflect.StructField) float64 {
	return getDefaultValue(field, reflect.Float64).(float64)
}

func defaultValueForKind(key string, kind reflect.Kind) any {
	switch kind {
	case reflect.String:
		return key
	case reflect.Bool:
		v, err := strconv.ParseBool(key)
		if err != nil {
			fmt.Printf("invalid default type %s for %v\n", key, kind.String())
			return false
		}
		return v
	case reflect.Int64:
		v, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			fmt.Printf("invalid default type %s for %v\n", key, kind.String())
			return int64(0)
		}
		return v
	case reflect.Float64:
		v, err := strconv.ParseFloat(key, 10)
		if err != nil {
			fmt.Printf("invalid default type %s for %v\n", key, kind.String())
			return float64(0)
		}
		return v
	default:
		return nil
	}
}

func getDefaultValueForKind(kind reflect.Kind) any {
	switch kind {
	case reflect.String:
		return ""
	case reflect.Bool:
		return false
	case reflect.Int64:
		return int64(0)
	case reflect.Float64:
		return float64(0)
	default:
		return nil
	}
}
