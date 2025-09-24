package models

// Users is a sample model you can modify or remove.
type Users struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email"    gorm:"uniqueIndex"`
	Username string `json:"username" gorm:"uniqueIndex"`
	Password string `json:"-"`
	Status   string `json:"status" gorm:"default:'active'"`
}

// Force GORM to use the singular table name "users".
func (Users) TableName() string { return "users" }
