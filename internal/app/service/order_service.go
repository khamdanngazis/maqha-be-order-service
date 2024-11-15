package service

import (
	"context"
	exEntity "maqhaa/order_service/external/entity"
	exRepo "maqhaa/order_service/external/repository"
	"maqhaa/order_service/internal/app/entity"
	"maqhaa/order_service/internal/app/model"
	"maqhaa/order_service/internal/app/repository"
	"sync"

	"github.com/go-playground/validator/v10"
)

type OrderService interface {
	AddOrder(context.Context, string, *model.OrderRequest) (*entity.Order, AppError)
	EditOrder(context.Context, string, *model.OrderRequest) (*entity.Order, AppError)
	GetOrder(context.Context, string, int) (*entity.Order, AppError)
	// Add more methods as needed
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo exRepo.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo exRepo.ProductRepository) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *orderService) AddOrder(ctx context.Context, token string, request *model.OrderRequest) (*entity.Order, AppError) {

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return nil, *NewInvalidRequestError(err.Error())
	}

	var orderDetails []entity.OrderDetail
	var totalPrice float64
	var wg sync.WaitGroup
	resultChan := make(chan *exEntity.Product, len(request.Orders))
	for _, v := range request.Orders {
		wg.Add(1)
		ordersDetail := v
		go func() {
			defer wg.Done()
			resultProduct, _ := s.productRepo.GetProductByID(ctx, ordersDetail.ProductID, token)
			resultChan <- resultProduct

		}()

	}

	wg.Wait()
	close(resultChan)
	for result := range resultChan {
		if result == nil {
			return nil, *NewProductNotFoundError()
		}
		for _, reqDetail := range request.Orders {
			if result.ID != reqDetail.ProductID {
				continue
			}
			if result.Price != reqDetail.Price {
				return nil, *NewInvalidProductPriceError()
			}
		}
	}

	for _, reqDetail := range request.Orders {
		totalPrice += reqDetail.Total

		// Convert the request detail to order detail entity
		orderDetail := entity.OrderDetail{
			ProductID: reqDetail.ProductID,
			Price:     reqDetail.Price,
			Quantity:  reqDetail.Quantity,
			Discount:  reqDetail.Discount,
			Total:     reqDetail.Total,
		}

		orderDetails = append(orderDetails, orderDetail)
	}

	if totalPrice != request.Total {
		return nil, *NewInvalidTotalError()
	}

	// Convert the request to the Order entity
	order := &entity.Order{
		ClientID:     uint(request.ClientID),
		CustomerName: request.CustomerName,
		PhoneNumber:  request.PhoneNumber,
		Total:        totalPrice,
		Status:       model.OrderStatusIncoming,
		OrderDetails: orderDetails,
		// Add other fields as needed
	}

	// Call the repository to add the order
	order, err := s.orderRepo.AddOrder(ctx, order)
	if err != nil {
		return nil, *NewUpdateQueryDBError()
	}

	return order, *NewSuccessError()
}

func (s *orderService) EditOrder(ctx context.Context, token string, request *model.OrderRequest) (*entity.Order, AppError) {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return nil, *NewInvalidRequestError(err.Error())
	}

	order, err := s.orderRepo.GetOrderByID(ctx, uint(request.ID), token)
	if err != nil {
		return nil, *NewOrderNotFoundError()
	}

	if order == nil {
		return nil, *NewOrderNotFoundError()
	}

	var totalPrice float64
	var orderDetails []entity.OrderDetail

	for _, reqDetail := range request.Orders {
		product, err := s.productRepo.GetProductByID(ctx, reqDetail.ProductID, token)
		if err != nil {
			return nil, *NewProductNotFoundError()
		}

		if product.Price != reqDetail.Price {
			return nil, *NewInvalidProductPriceError()
		}

		totalPrice += reqDetail.Total

		orderDetail := entity.OrderDetail{
			ProductID: reqDetail.ProductID,
			Price:     reqDetail.Price,
			Quantity:  reqDetail.Quantity,
			Discount:  reqDetail.Discount,
			Total:     reqDetail.Total,
		}

		orderDetails = append(orderDetails, orderDetail)
	}

	if totalPrice != request.Total {
		return nil, *NewInvalidTotalError()
	}

	order.Total = totalPrice
	order.OrderDetails = orderDetails
	order.CustomerName = request.CustomerName
	order.PhoneNumber = request.PhoneNumber
	order.Status = model.OrderStatusIncoming

	updatedOrder, err := s.orderRepo.EditOrder(ctx, order)
	if err != nil {
		return nil, *NewUpdateQueryDBError()
	}

	return updatedOrder, *NewSuccessError()
}

func (s orderService) GetOrder(ctx context.Context, token string, orderID int) (*entity.Order, AppError) {
	product, err := s.orderRepo.GetOrderByID(ctx, uint(orderID), token)
	if err != nil {
		return nil, *NewOrderNotFoundError()
	}

	if product == nil {
		return nil, *NewOrderNotFoundError()
	}
	return product, *NewSuccessError()
}
