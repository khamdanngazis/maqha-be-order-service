// auth_handler_test.go

package handler_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"maqhaa/library/helper"
	"maqhaa/library/logging"
	exEntity "maqhaa/order_service/external/entity"
	"maqhaa/order_service/internal/app/entity"
	"maqhaa/order_service/internal/app/model"
	"maqhaa/order_service/internal/app/repository"
	"maqhaa/order_service/internal/app/repository/mock"
	"maqhaa/order_service/internal/app/service"
	"maqhaa/order_service/internal/config"
	"maqhaa/order_service/internal/database"
	"maqhaa/order_service/internal/interface/http/handler"

	"gorm.io/gorm"
)

var db *gorm.DB
var orderHandler *handler.OrderHandler
var producRepo *mock.MockProductRepository

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	// Load testing configuration
	cfg, err := config.LoadConfig("../../cmd/config/config-test.yaml")
	if err != nil {
		panic(err)
	}

	// Set up a connection to the testing database
	db, err = database.NewDB(&cfg.Database)
	if err != nil {
		panic(err)
	}

	logFolder := flag.String("log.file", "../../logs", "Logging file")

	flag.Parse()

	// set logging file
	//logging.OutputScreen = true
	//logging.Filename = *logFile
	logging.InitLogger()
	currentDate := time.Now().Format("2006-01-02")
	logFilePath := fmt.Sprintf("%s/app_test_%s.log", *logFolder, currentDate)

	// Create the log file if it doesn't exist
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}

	// Set the logrus output to the log file
	logging.Log.SetOutput(logFile)

	// Apply database migrations (if any)
	// You can use db.AutoMigrate(&YourModel{}) to automatically apply migrations

	// Create a product service and handler
	producRepo = mock.NewMockProductRepository()
	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepository, producRepo)
	orderHandler = handler.NewOrderHandler(orderService)

}

func clearDB(tables []string) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	for _, v := range tables {
		sqlDB.Exec("delete from " + v)
	}

}

func tearDown() {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

}

func SampleUser() *entity.User {
	token, _ := helper.GenerateRandomString(16)
	return &entity.User{
		ID:           1,
		ClientID:     1,
		Username:     "sample",
		Password:     "rahasia",
		FullName:     "Sample User",
		Token:        token,
		TokenExpired: time.Now().Add(time.Hour), // Set a token expiration time
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
}

func SampleClient() *entity.Client {
	return &entity.Client{
		ID:          1,
		CompanyName: "Example Coffee",
		Email:       "info@examplecoffee.com",
		PhoneNumber: "+1234567890",
		Address:     "123 Main St, Cityville",
		OwnerName:   "John Doe",
		IsActive:    true,
		Token:       "JYA60sj03G6ii0LR3BfF", // Add a sample token if needed
		CreatedAt:   time.Now(),
	}
}

func SampleCategories(clientID uint) []*exEntity.ProductCategory {
	// Sample product data for Coffee category
	layoutFormat := "2006-01-02 15:04:05"
	coffeeProducts := []exEntity.Product{
		{
			ID:          1,
			Name:        "Espresso",
			Description: "Strong coffee",
			Price:       2.5,
			IsActive:    true,
			CreatedAt:   time.Now().Format(layoutFormat),
		},
		{
			ID:          2,
			Name:        "Latte",
			Description: "Coffee with milk",
			Price:       3.0,
			IsActive:    true,
			CreatedAt:   time.Now().Format(layoutFormat),
		},
	}

	// Sample product data for Tea category
	teaProducts := []exEntity.Product{
		{
			ID:          3,
			Name:        "Green Tea",
			Description: "Healthy tea",
			Price:       2.0,
			IsActive:    true,
			CreatedAt:   time.Now().Format(layoutFormat),
		},
		// Add more tea products as needed
	}

	// Sample product data for Snacks category
	snacksProducts := []exEntity.Product{
		{
			ID:          4,
			Name:        "Chips",
			Description: "Crispy snacks",
			Price:       1.5,
			IsActive:    true,
			CreatedAt:   time.Now().Format(layoutFormat),
		},
		// Add more snack products as needed
	}

	return []*exEntity.ProductCategory{
		{
			ID:        1,
			ClientID:  clientID,
			Name:      "Coffee",
			IsActive:  true,
			CreatedAt: time.Now(),
			Products:  coffeeProducts,
		},
		{
			ID:        2,
			ClientID:  clientID,
			Name:      "Tea",
			IsActive:  true,
			CreatedAt: time.Now(),
			Products:  teaProducts,
		},
		{
			ID:        3,
			ClientID:  clientID,
			Name:      "Snacks",
			IsActive:  true,
			CreatedAt: time.Now(),
			Products:  snacksProducts,
		},
		// Add more sample categories as needed
	}
}

func SampleOrder(clientID uint) *entity.Order {
	return &entity.Order{
		OrderNumber:  "ORD123",
		ClientID:     clientID,
		QueueNumber:  1,
		CustomerName: "John Doe",
		PhoneNumber:  "123456789",
		Total:        100.50,
		Status:       model.OrderStatusIncoming,
		StatusText:   "Pending",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UpdatedBy:    1,
		OrderDetails: []entity.OrderDetail{
			{
				ProductID: 1,
				Price:     50.25,
				Quantity:  2,
				Discount:  5.0,
				Total:     95.25,
			},
			{
				ProductID: 2,
				Price:     30.75,
				Quantity:  3,
				Discount:  2.5,
				Total:     88.25,
			},
		},
	}
}
