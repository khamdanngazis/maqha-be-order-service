package repository

import (
	"context"
	"fmt"
	"maqhaa/library/logging"
	"maqhaa/library/middleware"
	"maqhaa/order_service/internal/app/entity"
	"maqhaa/order_service/internal/app/model"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderRepository interface {
	AddOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	GetOrderByID(ctx context.Context, orderID uint, clientToken string) (*entity.Order, error)
	EditOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) AddOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	requestID, _ := ctx.Value(middleware.RequestIDKey).(string)
	tx := r.db.Begin()

	// Get the current date
	currentDate := time.Now().Format("2006-01-02")

	// Retrieve the latest queue number for the current date and client
	var latestQueueNumber int
	if err := tx.Table("order").
		Where("DATE(created_at) = ? AND client_id = ?", currentDate, order.ClientID).
		Select("IFNULL(MAX(queue_number), 0)").
		Set("gorm:query_option", "FOR UPDATE").
		Scan(&latestQueueNumber).
		Error; err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error AddOrder  %s", err.Error())
		tx.Rollback()
		return nil, err
	}

	// Increment the queue number
	latestQueueNumber++
	order.QueueNumber = latestQueueNumber

	// Generate the order number with leading zeros
	order.OrderNumber = fmt.Sprintf("ORD-%04d", latestQueueNumber)

	// Create the order
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error AddOrder  %s", err.Error())
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error AddOrder  %s", err.Error())
		return nil, err
	}

	return order, nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, orderID uint, clientToken string) (*entity.Order, error) {
	requestID, _ := ctx.Value(middleware.RequestIDKey).(string)
	var order entity.Order

	// Retrieve the order by ID and Client Token
	if err := r.db.Joins("JOIN client ON order.client_id = client.id").
		Where("order.id = ? AND client.token = ?", orderID, clientToken).
		First(&order).
		Error; err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error GetOrderByID  %s", err.Error())
		return nil, err
	}

	// Set the StatusText based on the Status value
	switch order.Status {
	case 1:
		order.StatusText = model.OrderStatusIncomingMessage
	case 2:
		order.StatusText = model.OrderStatusPaidMessage
	case 3:
		order.StatusText = model.OrderStatusProcessingMessage
	case 4:
		order.StatusText = model.OrderStatusSuccessMessage
	default:
		order.StatusText = "Unknown"
	}

	return &order, nil
}

func (r *orderRepository) EditOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	requestID, _ := ctx.Value(middleware.RequestIDKey).(string)
	tx := r.db.Begin()

	// Update order
	if err := tx.Model(&entity.Order{}).Where("id = ?", order.ID).Updates(order).Error; err != nil {
		tx.Rollback()
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error updating order %s", err.Error())
		return nil, err
	}

	// Remove all old order details
	if err := tx.Where("order_id = ?", order.ID).Delete(&entity.OrderDetail{}).Error; err != nil {
		tx.Rollback()
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error deleting old order details %s", err.Error())
		return nil, err
	}

	// Add new order details
	for _, detail := range order.OrderDetails {
		detail.OrderID = order.ID // Ensure the foreign key is set correctly
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error adding new order detail %s", err.Error())
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error committing transaction %s", err.Error())
		return nil, err
	}

	return order, nil
}
