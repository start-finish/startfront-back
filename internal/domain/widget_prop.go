package domain

type WidgetProp struct {
	ID        int     `db:"id" json:"id"`
	WidgetID  *int    `db:"widget_id" json:"widget_id"`
	PresetID  *int    `db:"preset_id" json:"preset_id"`
	Key       string  `db:"key" json:"key"`
	Value     string  `db:"value" json:"value"`
	ValueType string  `db:"value_type" json:"value_type"`
}
