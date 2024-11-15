package entity

import "time"

type Order struct {
	ID           uint          `gorm:"primary_key" json:"id"`
	OrderNumber  string        `json:"order_number"`
	ClientID     uint          `json:"client_id"`
	QueueNumber  int           `json:"queue_number"`
	CustomerName string        `json:"customer_name"`
	PhoneNumber  string        `json:"phone_number"`
	Total        float64       `json:"total"`
	Status       int           `json:"status"`
	StatusText   string        `json:"status_text" gorm:"-"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	UpdatedBy    int           `json:"updated_by"`
	OrderDetails []OrderDetail `json:"order_details,omitempty" gorm:"foreignkey:OrderID"`
}

func (Order) TableName() string {
	return "order"
}
