package domain

type ApplicationCollaborator struct {
	ID            int    `db:"id" json:"id"`
	ApplicationID int    `db:"application_id" json:"application_id"`
	UserID        int    `db:"user_id" json:"user_id"`
	Role          string `db:"role" json:"role"`
}
