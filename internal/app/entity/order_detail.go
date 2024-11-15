package entity

type OrderDetail struct {
	ID        uint    `gorm:"primary_key" json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Discount  float64 `json:"discount"`
	Total     float64 `json:"total"`
}

func (OrderDetail) TableName() string {
	return "order_detail"
}
