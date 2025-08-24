package models

// Admin_dashboard is a sample model you can modify or remove.
type Admin_dashboard struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

// Force GORM to use the singular table name "admin_dashboard".
func (Admin_dashboard) TableName() string { return "admin_dashboard" }
