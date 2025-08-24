package models

// Permissions is a sample model you can modify or remove.
type Permissions struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Code string `json:"code"`
	Description string `json:"description"`
}

// Force GORM to use the singular table name "permissions".
func (Permissions) TableName() string { return "permissions" }
