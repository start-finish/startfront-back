package pkg

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type BaseService[T any] struct {
	DB *gorm.DB
}

func (s *BaseService[T]) List() ([]T, error) {
	var items []T
	result := s.DB.Find(&items)
	return items, result.Error
}

func (s *BaseService[T]) GetByID(id uint) (*T, error) {
	var item T
	result := s.DB.First(&item, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (s *BaseService[T]) Create(item *T) error {
	return s.DB.Create(item).Error
}

func (s *BaseService[T]) Update(id uint, data map[string]interface{}) error {
	var existing T
	if err := s.DB.First(&existing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("record not found")
		}
		return err
	}

	return s.DB.Model(&existing).Updates(data).Error
}


func (s *BaseService[T]) Delete(id uint) error {
	var existing T
	if err := s.DB.First(&existing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("record not found")
		}
		return err
	}
	return s.DB.Delete(&existing).Error
}


// ✨ Duplicate checker
func (s *BaseService[T]) Exists(field string, model T) (bool, error) {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Capitalize first letter to match exported struct field
	fieldName := strings.ToUpper(field[:1]) + field[1:]
	fieldVal := v.FieldByName(fieldName)

	if !fieldVal.IsValid() {
		return false, fmt.Errorf("field %s not found", field)
	}

	var count int64
	if err := s.DB.Model(new(T)).Where(field+" = ?", fieldVal.Interface()).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
