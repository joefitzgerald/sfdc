package sfdc

import (
	"reflect"
	"strings"
)

// deepFields recurses struct types to fetch types embedded in a struct.
func deepFields(typ reflect.Type) []reflect.StructField {
	fields := make([]reflect.StructField, 0)
	for i := 0; i < typ.NumField(); i++ {
		v := typ.Field(i)
		switch v.Type.Kind() {
		case reflect.Struct:
			switch v.Name {
			case v.Type.Name():
				fields = append(fields, deepFields(v.Type)...)
			default:
				fields = append(fields, v)
			}
		default:
			fields = append(fields, v)
		}
	}

	return fields
}

func fieldsForType(t reflect.Type) string {
	result := []string{}
	fields := deepFields(t)
	for _, field := range fields {
		target := field.Name
		skip := false
		if jsonTag, ok := field.Tag.Lookup("json"); ok {
			jsonTag = strings.TrimSpace(jsonTag)
			if jsonTag == "-" {
				skip = true
				continue
			} else {
				segments := strings.Split(jsonTag, ",")
				if len(segments) > 0 && strings.TrimSpace(segments[0]) != "" {
					target = strings.TrimSpace(segments[0])
				}
			}
		}
		if sfdcTag, ok := field.Tag.Lookup("sfdc"); ok {
			sfdcTag = strings.TrimSpace(sfdcTag)
			if sfdcTag != "" {
				skip = false
				target = sfdcTag
			}
			if sfdcTag == "-" {
				skip = true
			}
		}
		if !skip {
			result = append(result, target)
		}
	}
	return strings.Join(result, ",")
}
