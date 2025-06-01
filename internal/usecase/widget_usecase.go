package usecase

import (
	"startfront-backend/internal/domain"
	"startfront-backend/internal/repository"
)

// CreateWidget creates a new widget
func CreateWidget(widget domain.Widget) error {
	return repository.InsertWidget(widget)
}

// // GetWidgetsByScreenID fetches widgets by screen ID
// func GetWidgetsByScreenID(screenID string) ([]domain.Widget, error) {
// 	return repository.GetWidgetsByScreenID(screenID)
// }

// // UpdateWidget updates widget details
// func UpdateWidget(id string, widget domain.Widget) error {
// 	return repository.UpdateWidget(id, widget)
// }

// // DeleteWidget deletes a widget
// func DeleteWidget(id string) error {
// 	return repository.DeleteWidget(id)
// }
