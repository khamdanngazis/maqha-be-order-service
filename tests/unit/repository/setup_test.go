package repository_test

import (
	"context"
	"flag"
	"fmt"
	"log"
	"maqhaa/library/logging"
	"maqhaa/library/middleware"
	exEntity "maqhaa/order_service/external/entity"
	"maqhaa/order_service/internal/app/entity"
	"maqhaa/order_service/internal/app/repository"
	"maqhaa/order_service/internal/app/repository/mock"
	"maqhaa/order_service/internal/config"
	"maqhaa/order_service/internal/database"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var db *gorm.DB
var orderRepo repository.OrderRepository
var producRepo *mock.MockProductRepository
var ctx context.Context

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	// Load testing configuration
	cfg, err := config.LoadConfig("../../../cmd/config/config-test.yaml")
	if err != nil {
		panic(err)
	}

	// Set up a connection to the testing database
	db, err = database.NewDB(&cfg.Database)
	if err != nil {
		panic(err)
	}

	logFolder := flag.String("log.file", "../../../logs", "Logging file")

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
	orderRepo = repository.NewOrderRepository(db)
	requestID := uuid.New().String()
	ctx = context.WithValue(context.Background(), middleware.RequestIDKey, requestID)
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
