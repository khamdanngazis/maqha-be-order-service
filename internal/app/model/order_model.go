package model

import "maqhaa/order_service/internal/app/entity"

const (
	OrderStatusIncoming          = 1
	OrderStatusIncomingMessage   = "Incoming"
	OrderStatusPaid              = 2
	OrderStatusPaidMessage       = "Paid"
	OrderStatusProcessing        = 3
	OrderStatusProcessingMessage = "Processing"
	OrderStatusSuccess           = 4
	OrderStatusSuccessMessage    = "Success"
)

type OrderRequest struct {
	ID           uint          `json:"order_id"`
	ClientID     uint          `json:"client_id" validate:"required"`
	CustomerName string        `json:"customer_name" validate:"required"`
	PhoneNumber  string        `json:"phone_number"`
	Total        float64       `json:"total" validate:"required,gt=0"`
	Orders       []OrderDetail `validate:"required,dive"`
}

type OrderDetail struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Price     float64 `json:"price" validate:"required,gt=0"`
	Quantity  int     `json:"quantity" validate:"required,gte=1"`
	Discount  float64 `json:"discount" validate:"gte=0"`
	Total     float64 `json:"total" validate:"required,gt=0"`
}

type OrderResponse struct {
	HTTPResponse
	Data *struct {
		OrderID     uint   `json:"order_id"`
		OrderNumber string `json:"order_number"`
	} `json:"data,omitempty"`
}

type GetOrderResponse struct {
	HTTPResponse
	Data *struct {
		Order *entity.Order `json:"order,omitempty"`
	} `json:"data,omitempty"`
}
