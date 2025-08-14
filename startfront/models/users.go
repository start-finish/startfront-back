package models

// Users is a sample model you can modify or remove.
type Users struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Force GORM to use the singular table name "users".
func (Users) TableName() string { return "users" }
