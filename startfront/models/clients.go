package models

// Clients is a sample model you can modify or remove.
type Clients struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Status string `json:"status" gorm:"default:'active'"`
}

// Force GORM to use the singular table name "clients".
func (Clients) TableName() string { return "clients" }
