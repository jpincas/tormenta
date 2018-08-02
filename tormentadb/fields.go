package tormentadb

import "reflect"

//MapFields returns a map keyed by fieldname with the value of the field as an interface
func MapFields(entity interface{}) map[string]interface{} {
	v := reflect.Indirect(reflect.ValueOf(entity))

	modelMap := map[string]interface{}{}
	fieldMap := map[string]interface{}{}

	for i := 0; i < v.NumField(); i++ {
		fieldType := v.Type().Field(i)
		fieldName := fieldType.Name

		// We want to flatten the embedded 'Model'
		// So recursively run this function on it,
		// and then merge the resulting map with the main 'fieldMap'
		if fieldName == "Model" {
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
func ListFields(entity Tormentable) []string {
	v := reflect.Indirect(reflect.ValueOf(entity))

	// We want ID, Created and LastUpdated to appear at the start
	// so we add those manually
	fields := []string{"ID", "Created", "LastUpdated"}

	// Run through all the fields, appending them to the list
	// Don't include 'Model'
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		if fieldName != "Model" {
			fields = append(fields, fieldName)
		}
	}

	return fields
}
