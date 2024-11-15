// product_handler_test.go

package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"maqhaa/library/logging"
	"maqhaa/library/middleware"
	"maqhaa/order_service/internal/app/entity"
	"maqhaa/order_service/internal/app/model"
	"maqhaa/order_service/internal/app/service"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestOrderProductHandler_Positive(t *testing.T) {
	tables := []string{"product", "product_category", "client", "user", "order_detail", "`order`"}
	defer clearDB(tables)

	client := SampleClient()

	categories := SampleCategories(client.ID)

	producRepo.SetProductResponse(categories[0].Products[0].ID, &categories[0].Products[0])
	producRepo.SetProductResponse(categories[0].Products[1].ID, &categories[0].Products[1])
	producRepo.SetProductResponse(categories[1].Products[0].ID, &categories[1].Products[0])
	producRepo.SetProductResponse(categories[2].Products[0].ID, &categories[2].Products[0])
	validRequest := model.OrderRequest{
		ClientID:     uint(client.ID),
		CustomerName: "John Doe",
		PhoneNumber:  "123456789",
		Total:        categories[0].Products[0].Price * 2,
		Orders: []model.OrderDetail{
			{ProductID: categories[0].Products[0].ID, Price: categories[0].Products[0].Price, Quantity: 2, Discount: 0.0, Total: categories[0].Products[0].Price * 2},
		},
	}

	orderRequestJSON, _ := json.Marshal(validRequest)
	req, err := http.NewRequest("POST", "/orders", bytes.NewBuffer(orderRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", client.Token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function

	http.HandlerFunc(orderHandler.CreateOrderHandler).ServeHTTP(rr, req)

	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response CreateOrderHandler")

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	var response model.OrderResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)
	assert.NotEmpty(t, response.Data.OrderNumber)
	assert.NotEmpty(t, response.Data.OrderID)

}

func TestOrderProductHandler_ProductNotFound(t *testing.T) {
	tables := []string{"product", "product_category", "client", "user", "order_detail", "`order`"}
	defer clearDB(tables)

	client := SampleClient()

	categories := SampleCategories(client.ID)

	producRepo.SetProductResponse(categories[0].Products[0].ID, &categories[0].Products[0])
	producRepo.SetProductResponse(categories[0].Products[1].ID, &categories[0].Products[1])
	producRepo.SetProductResponse(categories[1].Products[0].ID, &categories[1].Products[0])
	producRepo.SetProductResponse(categories[2].Products[0].ID, &categories[2].Products[0])
	validRequest := model.OrderRequest{
		ClientID:     uint(client.ID),
		CustomerName: "John Doe",
		PhoneNumber:  "123456789",
		Total:        categories[0].Products[0].Price * 2,
		Orders: []model.OrderDetail{
			{ProductID: 10, Price: categories[0].Products[0].Price, Quantity: 2, Discount: 0.0, Total: categories[0].Products[0].Price * 2},
		},
	}

	orderRequestJSON, _ := json.Marshal(validRequest)
	req, err := http.NewRequest("POST", "/orders", bytes.NewBuffer(orderRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", client.Token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function

	http.HandlerFunc(orderHandler.CreateOrderHandler).ServeHTTP(rr, req)

	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response CreateOrderHandler")

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	var response model.OrderResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, service.ProductNotFoundMessage, response.Message)
	assert.Equal(t, service.ProductNotFound, response.Code)
	assert.Nil(t, response.Data)

}

func TestOrderProductHandler_InvalidProductPrice(t *testing.T) {
	tables := []string{"product", "product_category", "client", "user", "order_detail", "`order`"}
	defer clearDB(tables)

	client := SampleClient()

	categories := SampleCategories(client.ID)

	producRepo.SetProductResponse(categories[0].Products[0].ID, &categories[0].Products[0])
	producRepo.SetProductResponse(categories[0].Products[1].ID, &categories[0].Products[1])
	producRepo.SetProductResponse(categories[1].Products[0].ID, &categories[1].Products[0])
	producRepo.SetProductResponse(categories[2].Products[0].ID, &categories[2].Products[0])
	validRequest := model.OrderRequest{
		ClientID:     uint(client.ID),
		CustomerName: "John Doe",
		PhoneNumber:  "123456789",
		Total:        categories[0].Products[0].Price * 2,
		Orders: []model.OrderDetail{
			{ProductID: categories[0].Products[0].ID, Price: categories[0].Products[1].Price, Quantity: 2, Discount: 0.0, Total: categories[0].Products[0].Price * 2},
		},
	}

	orderRequestJSON, _ := json.Marshal(validRequest)
	req, err := http.NewRequest("POST", "/orders", bytes.NewBuffer(orderRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", client.Token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function

	http.HandlerFunc(orderHandler.CreateOrderHandler).ServeHTTP(rr, req)

	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response CreateOrderHandler")

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	var response model.OrderResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, service.InvalidProductPriceMessage, response.Message)
	assert.Equal(t, service.InvalidProductPrice, response.Code)
	assert.Nil(t, response.Data)

}

func TestOrderProductHandler_InvalidTotal(t *testing.T) {
	tables := []string{"product", "product_category", "client", "user", "order_detail", "`order`"}
	defer clearDB(tables)

	client := SampleClient()

	categories := SampleCategories(client.ID)

	producRepo.SetProductResponse(categories[0].Products[0].ID, &categories[0].Products[0])
	producRepo.SetProductResponse(categories[0].Products[1].ID, &categories[0].Products[1])
	producRepo.SetProductResponse(categories[1].Products[0].ID, &categories[1].Products[0])
	producRepo.SetProductResponse(categories[2].Products[0].ID, &categories[2].Products[0])
	validRequest := model.OrderRequest{
		ClientID:     uint(client.ID),
		CustomerName: "John Doe",
		PhoneNumber:  "123456789",
		Total:        categories[0].Products[0].Price,
		Orders: []model.OrderDetail{
			{ProductID: categories[0].Products[0].ID, Price: categories[0].Products[0].Price, Quantity: 2, Discount: 0.0, Total: categories[0].Products[0].Price * 2},
		},
	}

	orderRequestJSON, _ := json.Marshal(validRequest)
	req, err := http.NewRequest("POST", "/orders", bytes.NewBuffer(orderRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", client.Token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function

	http.HandlerFunc(orderHandler.CreateOrderHandler).ServeHTTP(rr, req)

	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response CreateOrderHandler")

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	var response model.OrderResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, service.InvalidTotalMessage, response.Message)
	assert.Equal(t, service.InvalidTotal, response.Code)
	assert.Nil(t, response.Data)

}

func TestGetOrderHandler_Positive(t *testing.T) {
	tables := []string{"product", "product_category", "client", "user", "order_detail", "`order`"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(&client)
	order := SampleOrder(client.ID)
	errCreate := db.Create(&order).Error
	if errCreate != nil {
		t.Fatal(errCreate)
	}

	router := mux.NewRouter()
	router.HandleFunc("/order/{orderID}", orderHandler.GetOrderHandler).Methods("GET")

	req, err := http.NewRequest("GET", "/order/"+strconv.Itoa(int(order.ID)), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the client token in the request header
	req.Header.Set("Token", client.Token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	// Call the handler function
	router.ServeHTTP(rr, req)

	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response GetOrderHandler")

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	var response model.GetOrderResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)
	assert.NotEmpty(t, response.Data.Order.StatusText)
	assert.Equal(t, model.OrderStatusIncomingMessage, response.Data.Order.StatusText)
}

func TestGetOrderHandler_InvalidToken(t *testing.T) {
	tables := []string{"product", "product_category", "client", "user", "order_detail", "`order`"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(&client)
	order := SampleOrder(client.ID)
	errCreate := db.Create(&order).Error
	if errCreate != nil {
		t.Fatal(errCreate)
	}

	router := mux.NewRouter()
	router.HandleFunc("/order/{orderID}", orderHandler.GetOrderHandler).Methods("GET")

	req, err := http.NewRequest("GET", "/order/"+strconv.Itoa(int(order.ID)), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the client token in the request header
	req.Header.Set("Token", "")
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	// Call the handler function
	router.ServeHTTP(rr, req)

	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response GetOrderHandler")

	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response model.GetOrderResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, service.InvalidTokendMessage, response.Message)
	assert.Equal(t, service.InvalidToken, response.Code)
	assert.Nil(t, response.Data)
}

func TestGetOrderHandler_OrderNotFound(t *testing.T) {
	tables := []string{"product", "product_category", "client", "user", "order_detail", "`order`"}
	defer clearDB(tables)

	client := SampleClient()
	db.Create(&client)
	order := SampleOrder(client.ID)
	errCreate := db.Create(&order).Error
	if errCreate != nil {
		t.Fatal(errCreate)
	}

	router := mux.NewRouter()
	router.HandleFunc("/order/{orderID}", orderHandler.GetOrderHandler).Methods("GET")
	invalidOrderID := int(order.ID + 1)
	req, err := http.NewRequest("GET", "/order/"+strconv.Itoa(invalidOrderID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the client token in the request header
	req.Header.Set("Token", client.Token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	// Call the handler function
	router.ServeHTTP(rr, req)

	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response GetOrderHandler")

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	var response model.GetOrderResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, service.OrderNotFoundMessage, response.Message)
	assert.Equal(t, service.OrderNotFound, response.Code)
	assert.Nil(t, response.Data)
}

func TestEditProductHandler_Positive(t *testing.T) {
	tables := []string{"product", "product_category", "client", "user", "order_detail", "`order`"}
	defer clearDB(tables)

	client := SampleClient()

	categories := SampleCategories(client.ID)

	producRepo.SetProductResponse(categories[0].Products[0].ID, &categories[0].Products[0])
	producRepo.SetProductResponse(categories[0].Products[1].ID, &categories[0].Products[1])
	producRepo.SetProductResponse(categories[1].Products[0].ID, &categories[1].Products[0])
	producRepo.SetProductResponse(categories[2].Products[0].ID, &categories[2].Products[0])
	validRequest := model.OrderRequest{
		ClientID:     uint(client.ID),
		CustomerName: "John Doe",
		PhoneNumber:  "123456789",
		Total:        categories[0].Products[0].Price * 2,
		Orders: []model.OrderDetail{
			{ProductID: categories[0].Products[0].ID, Price: categories[0].Products[0].Price, Quantity: 2, Discount: 0.0, Total: categories[0].Products[0].Price * 2},
		},
	}

	orderRequestJSON, _ := json.Marshal(validRequest)
	req, err := http.NewRequest("POST", "/orders", bytes.NewBuffer(orderRequestJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Token", client.Token)
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the handler function

	http.HandlerFunc(orderHandler.CreateOrderHandler).ServeHTTP(rr, req)

	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response CreateOrderHandler")

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	var response model.OrderResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, service.SuccessMessage, response.Message)
	assert.Equal(t, service.SuccessError, response.Code)
	assert.NotEmpty(t, response.Data.OrderNumber)
	assert.NotEmpty(t, response.Data.OrderID)

	router := mux.NewRouter()
	router.HandleFunc("/orders/{orderID}", orderHandler.EditOrderHandler).Methods("PUT")

	validRequest = model.OrderRequest{
		ClientID:     uint(client.ID),
		CustomerName: "John Doe",
		PhoneNumber:  "123456789",
		Total:        (categories[0].Products[0].Price * 1) + (categories[0].Products[1].Price * 2),
		Orders: []model.OrderDetail{
			{ProductID: categories[0].Products[0].ID, Price: categories[0].Products[0].Price, Quantity: 1, Discount: 0.0, Total: categories[0].Products[0].Price * 1},
			{ProductID: categories[0].Products[1].ID, Price: categories[0].Products[1].Price, Quantity: 2, Discount: 0.0, Total: categories[1].Products[0].Price * 2},
		},
	}

	orderRequestJSON, _ = json.Marshal(validRequest)

	req, err = http.NewRequest("PUT", "/orders/"+strconv.Itoa(int(response.Data.OrderID)), bytes.NewBuffer(orderRequestJSON))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Token", client.Token)
	requestID = uuid.New().String()
	ctx = context.WithValue(req.Context(), middleware.RequestIDKey, requestID)
	req = req.WithContext(ctx)
	// Create a response recorder to capture the handler's response
	rr = httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	logging.Log.WithFields(logrus.Fields{
		"RequestID": requestID,
		"Status":    rr.Code,
		"Body":      rr.Body.String(),
	}).Info("Outgoing response")
	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	var orders entity.Order
	err = db.Preload("OrderDetails").Where("id = ?", response.Data.OrderID).First(&orders).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, orders)
	assert.Equal(t, 2, len(orders.OrderDetails))
	assert.Equal(t, validRequest.Orders[0].Quantity, orders.OrderDetails[0].Quantity)
	assert.Equal(t, validRequest.Orders[1].Quantity, orders.OrderDetails[1].Quantity)

}
