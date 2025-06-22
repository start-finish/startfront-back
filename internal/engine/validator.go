package engine

import (
	"errors"
)

func ValidateData(schema Schema, data map[string]interface{}) error {
	for _, field := range schema.Fields {
		if field.Required {
			if value, ok := data[field.Name]; !ok || value == "" {
				return errors.New("missing required field: " + field.Name)
			}
		}
	}
	return nil
}
