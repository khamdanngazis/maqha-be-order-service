package repository_test

import (
	"maqhaa/order_service/internal/app/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrderRepository_AddOrder(t *testing.T) {
	tables := []string{"order_detail", "`order`"}
	defer clearDB(tables)

	// Test data
	order := &entity.Order{
		ClientID:     1,
		CustomerName: "John Doe",
		PhoneNumber:  "123456789",
		Total:        100.0,
		Status:       1,
		OrderDetails: []entity.OrderDetail{
			{
				ProductID: 1,
				Price:     50.0,
				Quantity:  2,
				Discount:  10.0,
				Total:     90.0,
			},
		},
	}

	// Test AddOrder function
	createdOrder, err := orderRepo.AddOrder(ctx, order)
	assert.NoError(t, err)
	assert.NotNil(t, createdOrder)
	assert.NotEmpty(t, createdOrder.ID)
	assert.NotEmpty(t, createdOrder.OrderNumber)
	assert.Equal(t, order.ClientID, createdOrder.ClientID)
	assert.Equal(t, order.CustomerName, createdOrder.CustomerName)
	assert.Equal(t, order.PhoneNumber, createdOrder.PhoneNumber)
	assert.Equal(t, order.Total, createdOrder.Total)
	assert.Equal(t, order.Status, createdOrder.Status)
	assert.Equal(t, "ORD-0001", order.OrderNumber)
	assert.Equal(t, 1, order.QueueNumber)

	// Verify the order details
	assert.Len(t, createdOrder.OrderDetails, 1)
	assert.NotEmpty(t, createdOrder.OrderDetails[0].ID)
	assert.Equal(t, order.OrderDetails[0].ProductID, createdOrder.OrderDetails[0].ProductID)
	assert.Equal(t, order.OrderDetails[0].Price, createdOrder.OrderDetails[0].Price)
	assert.Equal(t, order.OrderDetails[0].Quantity, createdOrder.OrderDetails[0].Quantity)
	assert.Equal(t, order.OrderDetails[0].Discount, createdOrder.OrderDetails[0].Discount)
	assert.Equal(t, order.OrderDetails[0].Total, createdOrder.OrderDetails[0].Total)
}

func TestOrderRepository_AddOrderDouble(t *testing.T) {
	tables := []string{"order_detail", "`order`"}
	defer clearDB(tables)

	// Create two orders with order details
	order1 := &entity.Order{
		ClientID:     1,
		CustomerName: "John Doe",
		PhoneNumber:  "123456789",
		Total:        100.0,
		Status:       1,
		CreatedAt:    time.Now(),
		OrderDetails: []entity.OrderDetail{
			{ProductID: 1, Price: 50.0, Quantity: 2, Discount: 5.0},
		},
	}

	order2 := &entity.Order{
		ClientID:     1,
		CustomerName: "Jane Doe",
		PhoneNumber:  "987654321",
		Total:        150.0,
		Status:       1,
		CreatedAt:    time.Now(),
		OrderDetails: []entity.OrderDetail{
			{ProductID: 2, Price: 75.0, Quantity: 2, Discount: 7.5},
		},
	}

	// Add the first order
	resultOrder1, err := orderRepo.AddOrder(ctx, order1)
	assert.NoError(t, err)
	assert.NotNil(t, resultOrder1)
	assert.NotEmpty(t, resultOrder1.OrderNumber)
	assert.Equal(t, 1, resultOrder1.QueueNumber)
	assert.Equal(t, "ORD-0001", resultOrder1.OrderNumber)
	// Add the second order
	resultOrder2, err := orderRepo.AddOrder(ctx, order2)
	assert.NoError(t, err)
	assert.NotNil(t, resultOrder2)
	assert.NotEmpty(t, resultOrder2.OrderNumber)
	assert.Equal(t, 2, resultOrder2.QueueNumber)
	assert.Equal(t, "ORD-0002", resultOrder2.OrderNumber)

}

func TestOrderRepository_AddOrderDoubleClient(t *testing.T) {
	tables := []string{"order_detail", "`order`"}
	defer clearDB(tables)

	// Create two orders with order details
	order1 := &entity.Order{
		ClientID:     1,
		CustomerName: "John Doe",
		PhoneNumber:  "123456789",
		Total:        100.0,
		Status:       1,
		CreatedAt:    time.Now(),
		OrderDetails: []entity.OrderDetail{
			{ProductID: 1, Price: 50.0, Quantity: 2, Discount: 5.0},
		},
	}

	order2 := &entity.Order{
		ClientID:     2,
		CustomerName: "Jane Doe",
		PhoneNumber:  "987654321",
		Total:        150.0,
		Status:       1,
		CreatedAt:    time.Now(),
		OrderDetails: []entity.OrderDetail{
			{ProductID: 2, Price: 75.0, Quantity: 2, Discount: 7.5},
		},
	}

	// Add the first order
	resultOrder1, err := orderRepo.AddOrder(ctx, order1)
	assert.NoError(t, err)
	assert.NotNil(t, resultOrder1)
	assert.NotEmpty(t, resultOrder1.OrderNumber)
	assert.Equal(t, 1, resultOrder1.QueueNumber)
	assert.Equal(t, "ORD-0001", resultOrder1.OrderNumber)

	// Add the second order
	resultOrder2, err := orderRepo.AddOrder(ctx, order2)
	assert.NoError(t, err)
	assert.NotNil(t, resultOrder2)
	assert.NotEmpty(t, resultOrder2.OrderNumber)
	assert.Equal(t, 1, resultOrder2.QueueNumber)
	assert.Equal(t, "ORD-0001", resultOrder2.OrderNumber)
}

func TestOrderRepository_EditOrderDoubleClient(t *testing.T) {
	tables := []string{"order_detail", "`order`"}
	defer clearDB(tables)

	// Create two orders with order details
	order := &entity.Order{
		ClientID:     1,
		CustomerName: "John Doe",
		PhoneNumber:  "123456789",
		Total:        100.0,
		Status:       1,
		CreatedAt:    time.Now(),
		OrderDetails: []entity.OrderDetail{
			{ProductID: 1, Price: 50.0, Quantity: 2, Discount: 5.0, Total: 100.0},
		},
	}

	// Add the first order
	resultOrder1, err := orderRepo.AddOrder(ctx, order)
	assert.NoError(t, err)
	assert.NotNil(t, resultOrder1)
	assert.NotEmpty(t, resultOrder1.OrderNumber)
	assert.Equal(t, 1, resultOrder1.QueueNumber)
	assert.Equal(t, "ORD-0001", resultOrder1.OrderNumber)

	// Query table order menggunakan db GORM dengan kondisi order i1 = 1
	newOrder := &entity.Order{
		ID:           resultOrder1.ID,
		OrderNumber:  resultOrder1.OrderNumber,
		ClientID:     1,
		CustomerName: "John Doe",
		PhoneNumber:  "123456789",
		Total:        200.0,
		Status:       1,
		CreatedAt:    time.Now(),
		OrderDetails: []entity.OrderDetail{
			{ProductID: 1, Price: 50.0, Quantity: 1, Discount: 5.0, Total: 50.0},
			{ProductID: 2, Price: 75.0, Quantity: 2, Discount: 7.5, Total: 150.0},
		},
	}

	resultEdit, err := orderRepo.EditOrder(ctx, newOrder)
	assert.NoError(t, err)

	var orders entity.Order
	err = db.Preload("OrderDetails").Where("id = ?", resultEdit.ID).First(&orders).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, orders)
	assert.Equal(t, 2, len(orders.OrderDetails))
	assert.Equal(t, newOrder.OrderDetails[0].Quantity, orders.OrderDetails[0].Quantity)
	assert.Equal(t, newOrder.OrderDetails[1].Quantity, orders.OrderDetails[1].Quantity)

}
