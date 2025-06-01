package domain

type Widget struct {
	ID         int    `db:"id" json:"id"`
	ScreenID   int    `db:"screen_id" json:"screen_id"`
	ParentID   *int   `db:"parent_id" json:"parent_id"`
	PresetID   *int   `db:"preset_id" json:"preset_id"`
	Type       string `db:"type" json:"type"`
	ChildIndex *int   `db:"child_index" json:"child_index"`
	AuthBy     *int   `db:"auth_by" json:"auth_by"`
}
