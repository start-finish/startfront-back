package domain

type Screen struct {
	ID          int    `db:"id" json:"id"`
	AppID       int    `db:"application_id" json:"application_id"`
	Name        string `db:"name" json:"name"`
	Code        string `db:"code" json:"code"`
	Route       string `db:"route" json:"route"`
	Description string `db:"description" json:"description"`
	Params      string `db:"params" json:"params"`
	Validate    string `db:"validate" json:"validate"`
	AuthBy      *int   `db:"auth_by" json:"auth_by"`
	AuthByName  string `json:"auth_name"`
}
