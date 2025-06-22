package engine

import (
	"errors"
	"reflect"
	"time"
)

func BuildGormModel(schemaDef Schema) (interface{}, error) {
	fields := make([]reflect.StructField, 0)

	for _, f := range schemaDef.Fields {
		goType, err := mapFieldType(f.Type)
		if err != nil {
			return nil, err
		}

		tag := buildGormTag(f)

		fields = append(fields, reflect.StructField{
			Name: capitalize(f.Name),
			Type: goType,
			Tag:  reflect.StructTag(`gorm:"` + tag + `" json:"` + f.Name + `"`),
		})
	}

	structType := reflect.StructOf(fields)
	modelPtr := reflect.New(structType).Interface()

	return modelPtr, nil
}

func mapFieldType(fieldType string) (reflect.Type, error) {
	switch fieldType {
	case "string":
		return reflect.TypeOf(""), nil
	case "uint":
		return reflect.TypeOf(uint(0)), nil
	case "int":
		return reflect.TypeOf(int(0)), nil
	case "float":
		return reflect.TypeOf(float64(0)), nil
	case "bool":
		return reflect.TypeOf(false), nil
	case "timestamp":
		return reflect.TypeOf(time.Time{}), nil
	case "jsonb":
		return reflect.TypeOf(map[string]interface{}{}), nil
	default:
		return nil, errors.New("unsupported field type: " + fieldType)
	}
}

func buildGormTag(f Field) string {
	tag := ""

	if f.PrimaryKey {
		tag += "primaryKey;"
	}
	if f.AutoIncrement {
		tag += "autoIncrement;"
	}
	if f.Required {
		tag += "not null;"
	}
	if f.Unique {
		tag += "unique;"
	}
	if f.Index {
		tag += "index;"
	}

	return tag
}

func capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	return string(s[0]-32) + s[1:]
}

func SetModelFields(model interface{}, data map[string]interface{}) {
	v := reflect.ValueOf(model).Elem()

	for key, value := range data {
		field := v.FieldByName(capitalize(key))
		if field.IsValid() && field.CanSet() {
			targetType := field.Type()

			val := reflect.ValueOf(value)

			// Convert float64 → int/uint
			if val.Type().Kind() == reflect.Float64 {
				switch targetType.Kind() {
				case reflect.Int:
					val = reflect.ValueOf(int(val.Float()))
				case reflect.Int64:
					val = reflect.ValueOf(int64(val.Float()))
				case reflect.Uint:
					val = reflect.ValueOf(uint(val.Float()))
				case reflect.Uint64:
					val = reflect.ValueOf(uint64(val.Float()))
				}
			}

			// If convertible, set field
			if val.Type().ConvertibleTo(targetType) {
				field.Set(val.Convert(targetType))
			}
		}
	}
}
