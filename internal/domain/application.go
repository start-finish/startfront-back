package domain

type Application struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Code        string `db:"code" json:"code"`
	Route       string `db:"route" json:"route"`
	Description string `db:"description" json:"description"`
	AuthBy      *int   `db:"auth_by" json:"auth_by"`
	AuthByName  string `json:"auth_name"`
}
