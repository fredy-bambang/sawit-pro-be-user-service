package models

// User model
type User struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phone" gorm:"unique,not null"`
	Fullname    string `json:"username" gorm:"not null"`
	Password    string `json:"password" gorm:"not null"`
	SaltToken   string `json:"salt_token" gorm:"not null"`
}
