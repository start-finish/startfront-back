package repository

import (
	"startfront-backend/internal/domain"
)

// InsertWidget creates a new widget in the database
func InsertWidget(widget domain.Widget) error {
	// _, err := db.DB.Exec("INSERT INTO widgets (screen_id, type, props) VALUES (?, ?, ?)", widget.ScreenID, widget.Type, widget.Props)
	// return err
	return nil
}

// GetWidgetsByScreenID fetches widgets by screen ID
// func GetWidgetsByScreenID(screenID string) ([]domain.Widget, error) {
// 	// rows, err := db.DB.Query("SELECT id, screen_id, type, props FROM widgets WHERE screen_id = ?", screenID)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// defer rows.Close()

// 	// var widgets []domain.Widget
// 	// for rows.Next() {
// 	// 	var widget domain.Widget
// 	// 	if err := rows.Scan(&widget.ID, &widget.ScreenID, &widget.Type, &widget.Props); err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// 	widgets = append(widgets, widget)
// 	// }
// 	// return widgets, nil
// 	return nil
// }

// UpdateWidget updates a widget
func UpdateWidget(id string, widget domain.Widget) error {
	// _, err := db.DB.Exec("UPDATE widgets SET type = ?, props = ? WHERE id = ?", widget.Type, widget.Props, id)
	// return err
	return nil
}

// DeleteWidget deletes a widget by ID
func DeleteWidget(id string) error {
	// _, err := db.DB.Exec("DELETE FROM widgets WHERE id = ?", id)
	// return err
	return nil
}
