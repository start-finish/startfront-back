package models

// Users is a sample model you can modify or remove.
type Users struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Email       string `json:"email"    gorm:"uniqueIndex"`
	Username    string `json:"username" gorm:"uniqueIndex"`
	Password    string `json:"-"`
	Status      string `json:"status" gorm:"default:'active'"`
	Role        string `json:"role"`
	Description string `json:"description"`
}

// Force GORM to use the singular table name "users".
func (Users) TableName() string { return "users" }
