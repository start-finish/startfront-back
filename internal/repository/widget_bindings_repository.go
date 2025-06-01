package repository

import "startfront-backend/internal/domain"

func InsertWidgetBinding(wb domain.WidgetBinding) error {
	query := `
	INSERT INTO widget_bindings 
	(widget_id, connection_id, endpoint, method, query_params, headers, body, response_path, binding_key)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := db.Exec(query,
		wb.WidgetID,
		wb.ConnectionID,
		wb.Endpoint,
		wb.Method,
		wb.QueryParams,
		wb.Headers,
		wb.Body,
		wb.ResponsePath,
		wb.BindingKey,
	)
	return err
}

func GetBindingsByWidgetID(widgetID int) ([]domain.WidgetBinding, error) {
	var bindings []domain.WidgetBinding
	err := db.Select(&bindings, `SELECT * FROM widget_bindings WHERE widget_id = $1`, widgetID)
	return bindings, err
}
