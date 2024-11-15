package entity

import (
	"time"
)

type ProductCategory struct {
	ID        uint      `json:"id"`
	ClientID  uint      `json:"clientId"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	Products  []Product `json:"products"`
}
