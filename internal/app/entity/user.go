package entity

import "time"

// User represents a user in the system.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ClientID     uint      `json:"clientId"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	FullName     string    `json:"fullName"`
	Token        string    `json:"token"`
	TokenExpired time.Time `json:"tokenExpired"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (User) TableName() string {
	return "user"
}
