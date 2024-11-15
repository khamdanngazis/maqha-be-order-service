package entity

type Product struct {
	ID          uint    `json:"id"`
	CategoryID  uint    `json:"categoryId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	IsActive    bool    `json:"isActive"`
	CreatedAt   string  `json:"createdAt"`
}
