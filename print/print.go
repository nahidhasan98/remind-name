package console

import (
	"fmt"
	"reflect"
	"strings"
)

var defaultIndent = 0

func PrettyPrint(v interface{}, indent int) {
	if defaultIndent == 0 {
		defaultIndent = indent
	}

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	switch val.Kind() {
	case reflect.Struct:
		// Print struct fields
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			fieldValue := val.Field(i)

			// If the field is a struct, map, or pointer, recursively print its fields
			if fieldValue.Kind() == reflect.Struct || fieldValue.Kind() == reflect.Map || fieldValue.Kind() == reflect.Ptr {
				fmt.Printf("%s%s:\n", strings.Repeat(" ", indent), field.Name)
				PrettyPrint(fieldValue.Interface(), indent+defaultIndent)
			} else {
				// Otherwise, print the key and value
				fmt.Printf("%s%s: %v\n", strings.Repeat(" ", indent), field.Name, fieldValue.Interface())
			}
		}
	case reflect.Map:
		// Print map key-value pairs
		for _, key := range val.MapKeys() {
			mapValue := val.MapIndex(key)

			// If the value is a struct, map, or pointer, recursively print its fields
			if mapValue.Kind() == reflect.Struct || mapValue.Kind() == reflect.Map || mapValue.Kind() == reflect.Ptr {
				fmt.Printf("%s%v:\n", strings.Repeat(" ", indent), key)
				PrettyPrint(mapValue.Interface(), indent+defaultIndent)
			} else {
				// Otherwise, print the key and value
				fmt.Printf("%s%v: %v\n", strings.Repeat(" ", indent), key, mapValue.Interface())
			}
		}
	case reflect.Slice, reflect.Array:
		// Print slice or array elements
		fmt.Printf("%s[\n", strings.Repeat(" ", indent))
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			PrettyPrint(elem.Interface(), indent+defaultIndent)
		}
		fmt.Printf("%s]\n", strings.Repeat(" ", indent))
	case reflect.Ptr:
		// Print pointer value (dereference)
		if !val.IsNil() {
			fmt.Printf("%s*Pointer:\n", strings.Repeat(" ", indent))
			PrettyPrint(val.Elem().Interface(), indent+defaultIndent)
		} else {
			fmt.Printf("%s<nil pointer>\n", strings.Repeat(" ", indent))
		}
	case reflect.String, reflect.Int, reflect.Float64, reflect.Bool:
		// Print basic types
		fmt.Printf("%s%v\n", strings.Repeat(" ", indent), val.Interface())
	default:
		fmt.Printf("%s<unsupported type: %s>\n", strings.Repeat(" ", indent), typ.Kind())
	}
}
