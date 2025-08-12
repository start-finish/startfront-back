package models

// Login is a sample model you can modify or remove.
type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}
