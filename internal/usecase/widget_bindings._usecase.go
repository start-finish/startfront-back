package usecase

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/repository"
)

func CreateWidgetBinding(binding domain.WidgetBinding) error {
	return repository.InsertWidgetBinding(binding)
}

func GetWidgetBindings(widgetID int) ([]domain.WidgetBinding, error) {
	return repository.GetBindingsByWidgetID(widgetID)
}
