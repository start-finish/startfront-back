package domain

type AppConnection struct {
	ID            int    `db:"id" json:"id"`
	ApplicationID int    `db:"application_id" json:"application_id"`
	Name          string `db:"name" json:"name"`
	Type          string `db:"type" json:"type"`
	Config        string `db:"config" json:"config"`
}
