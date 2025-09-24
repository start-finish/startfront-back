package models

// Projects is a sample model you can modify or remove.
type Projects struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	ClientId int64  `json:"client_id"`
	Name     string `json:"name"`
	Status   string `json:"status" gorm:"default:'active'"`

	Client Clients `gorm:"foreignKey:ClientId;onDelete:CASCADE" json:"clients"`
}

// Force GORM to use the singular table name "projects".
func (Projects) TableName() string { return "projects" }
