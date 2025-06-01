package domain

type WidgetBinding struct {
	ID           int     `db:"id" json:"id"`
	WidgetID     int     `db:"widget_id" json:"widget_id"`
	ConnectionID int     `db:"connection_id" json:"connection_id"`
	Endpoint     string  `db:"endpoint" json:"endpoint"`
	Method       string  `db:"method" json:"method"`
	QueryParams  *string `db:"query_params" json:"query_params"`
	Headers      *string `db:"headers" json:"headers"`
	Body         *string `db:"body" json:"body"`
	ResponsePath *string `db:"response_path" json:"response_path"`
	BindingKey   *string `db:"binding_key" json:"binding_key"`
}
