package types

import (
	"log/slog"
	"reflect"
)

func wrapUnion(value interface{}, typeName string) map[string]interface{} {
	slog.Debug("Wrapping union", "value", value, "type", typeName)
	if value == nil {
		return map[string]interface{}{"null": nil}
	}

	val := reflect.ValueOf(value)

	// Handle cases where the value can be nil (pointers, slices, maps, interfaces, channels)
	if val.Kind() == reflect.Ptr || val.Kind() == reflect.Slice || val.Kind() == reflect.Map || val.Kind() == reflect.Interface || val.Kind() == reflect.Chan {
		if val.IsNil() {
			return map[string]interface{}{"null": nil}
		}
		// Dereference pointer if not nil for further handling
		if val.Kind() == reflect.Ptr {
			value = val.Elem().Interface()
		}
	}

	// Handle cases where the value can be nil or empty (pointers, slices, maps, interfaces, channels)
	if val.Kind() == reflect.Slice && val.Len() == 0 {
		// For an empty slice or nil value, return the appropriate Avro schema representation for an empty array
		return map[string]interface{}{"array": []interface{}{}}
	}

	switch typeName {
	default:
		slog.Debug("Wrapping default", "value", value, "type", typeName)
		return map[string]interface{}{typeName: value}
	}
}
