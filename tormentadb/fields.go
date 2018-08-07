package tormentadb

import (
	"reflect"
)

//MapFields returns a map keyed by fieldname with the value of the field as an interface
func MapFields(entity interface{}) map[string]interface{} {
	v := reflect.Indirect(reflect.ValueOf(entity))

	modelMap := map[string]interface{}{}
	fieldMap := map[string]interface{}{}

	for i := 0; i < v.NumField(); i++ {
		fieldType := v.Type().Field(i)
		fieldName := fieldType.Name

		// Recursively flatten embedded structs
		if fieldType.Type.Kind() == reflect.Struct && fieldType.Anonymous {
			modelMap = MapFields(v.Field(i).Interface())
			for k, v := range modelMap {
				fieldMap[k] = v
			}
		} else {
			fieldMap[fieldName] = v.Field(i).Interface()
		}
	}

	return fieldMap
}

// ListFields returns a list of fields for the entity, with ID, Created and LastUpdated always at the start
func ListFields(entity interface{}) []string {
	// We want ID, Created and LastUpdated to appear at the start
	// so we add those manually
	return append([]string{"ID", "Created", "LastUpdated"}, structFields(entity)...)
}

func structFields(entity interface{}) (fields []string) {
	v := reflect.Indirect(reflect.ValueOf(entity))

	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		fieldType := v.Type().Field(i)

		// Recursively flatten embedded structs - don't include 'Model'
		if fieldType.Type.Kind() == reflect.Struct &&
			fieldType.Anonymous &&
			fieldName != "Model" {
			l := structFields(v.Field(i).Interface())
			fields = append(fields, l...)
		} else if fieldName != "Model" {
			fields = append(fields, fieldName)
		}
	}

	return
}
