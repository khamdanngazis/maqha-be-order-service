// internal/handler/Order_handler.go

package handler

import (
	"encoding/json"
	"maqhaa/library/logging"
	"maqhaa/library/middleware"
	"maqhaa/order_service/internal/app/entity"
	"maqhaa/order_service/internal/app/model"
	"maqhaa/order_service/internal/app/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// OrderHandler handles HTTP requests related to user Orderentication and Orderorization.
type OrderHandler struct {
	orderService service.OrderService
}

// NewOrderHandler creates a new OrderHandler instance.
func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// LoginHandler handles the HTTP request for user login.
func (h *OrderHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var orderRequest model.OrderRequest
	var orderResponse model.OrderResponse
	var appError service.AppError

	logID, _ := r.Context().Value(middleware.RequestIDKey).(string)
	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid request payload")

		appError = *service.NewInvalidFormatError()
		orderResponse = model.OrderResponse{
			HTTPResponse: *model.NewHTTPResponse(appError.Code, appError.Message, nil),
		}
		sendJSONResponse(w, orderResponse, appError.Code)
		return
	}

	// Call the order service to add the order
	order, appErr := h.orderService.AddOrder(r.Context(), token, &orderRequest)

	orderResponse = model.OrderResponse{
		HTTPResponse: *model.NewHTTPResponse(appErr.Code, appErr.Message, nil),
	}
	if appErr.Code != service.SuccessError {
		// Handle application-specific errors
		orderResponse.Data = nil
		sendJSONResponse(w, orderResponse, appError.Code)
		return
	}

	orderResponse.Data = &struct {
		OrderID     uint   `json:"order_id"`
		OrderNumber string `json:"order_number"`
	}{
		OrderID:     uint(order.ID),
		OrderNumber: order.OrderNumber,
	}

	sendJSONResponse(w, orderResponse, appError.Code)
}

func (h *OrderHandler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	var orderResponse model.GetOrderResponse
	var appError service.AppError

	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["orderID"])

	if err != nil {
		appError = *service.NewOrderNotFoundError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	// Call the order service to add the order
	order, appErr := h.orderService.GetOrder(r.Context(), token, orderID)

	orderResponse = model.GetOrderResponse{
		HTTPResponse: *model.NewHTTPResponse(appErr.Code, appErr.Message, nil),
	}
	if appErr.Code != service.SuccessError {
		// Handle application-specific errors
		orderResponse.Data = nil
		sendJSONResponse(w, orderResponse, appError.Code)
		return
	}

	orderResponse.Data = &struct {
		Order *entity.Order `json:"order,omitempty"`
	}{
		Order: order,
	}

	sendJSONResponse(w, orderResponse, appError.Code)
}

// EditOrderHandler handles the HTTP request for editing an order.
func (h *OrderHandler) EditOrderHandler(w http.ResponseWriter, r *http.Request) {
	var orderRequest model.OrderRequest
	var orderResponse model.OrderResponse
	var appError service.AppError

	logID, _ := r.Context().Value(middleware.RequestIDKey).(string)
	token := r.Header.Get("Token")

	if token == "" {
		appError = *service.NewInvalidTokenError()
		response := model.NewHTTPResponse(appError.Code, appError.Message, nil)
		sendJSONResponse(w, response, appError.Code)
		return
	}

	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["orderID"])
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid order ID format")
		appError = *service.NewOrderNotFoundError()
		orderResponse = model.OrderResponse{
			HTTPResponse: *model.NewHTTPResponse(appError.Code, appError.Message, nil),
		}
		sendJSONResponse(w, orderResponse, appError.Code)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": logID}).Info("Invalid request payload")
		appError = *service.NewInvalidFormatError()
		orderResponse = model.OrderResponse{
			HTTPResponse: *model.NewHTTPResponse(appError.Code, appError.Message, nil),
		}
		sendJSONResponse(w, orderResponse, appError.Code)
		return
	}

	orderRequest.ID = uint(orderID)

	// Call the order service to update the order
	order, appErr := h.orderService.EditOrder(r.Context(), token, &orderRequest)

	orderResponse = model.OrderResponse{
		HTTPResponse: *model.NewHTTPResponse(appErr.Code, appErr.Message, nil),
	}
	if appErr.Code != service.SuccessError {
		// Handle application-specific errors
		orderResponse.Data = nil
		sendJSONResponse(w, orderResponse, appError.Code)
		return
	}

	orderResponse.Data = &struct {
		OrderID     uint   `json:"order_id"`
		OrderNumber string `json:"order_number"`
	}{
		OrderID:     uint(order.ID),
		OrderNumber: order.OrderNumber,
	}

	sendJSONResponse(w, orderResponse, appErr.Code)
}
