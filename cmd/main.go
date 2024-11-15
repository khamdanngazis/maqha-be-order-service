// cmd/main.go

package main

import (
	"flag"
	"fmt"
	"log"
	"maqhaa/library/logging"
	exRepo "maqhaa/order_service/external/repository"
	"maqhaa/order_service/internal/app/repository"
	"maqhaa/order_service/internal/app/service"
	"maqhaa/order_service/internal/config"
	"maqhaa/order_service/internal/database"
	"maqhaa/order_service/internal/interface/http/handler"
	"maqhaa/order_service/internal/interface/http/router"
	"os"
	"time"
)

func main() {
	// Define a command line flag for the config file path
	configFilePath := flag.String("config", "config/config.yaml", "path to the config file")
	logFile := flag.String("log.file", "../logs", "Logging file")

	flag.Parse()

	initLogging(*logFile)

	// Load the configuration
	cfg, err := config.LoadConfig(*configFilePath)
	if err != nil {
		logging.Log.Fatalf("Error loading configuration: %v", err)
	}
	logging.Log.Infof("Load configuration from %v", *configFilePath)
	// Access configuration values
	dbConfig := cfg.Database

	db, err := database.NewDB(&dbConfig)
	if err != nil {
		logging.Log.Fatalf("Error loading configuration: %v", err)
	}

	// Close the database connection when done
	sqlDB, err := db.DB()
	if err != nil {
		logging.Log.Fatalf("Error getting DB connection: %v", err)
	}
	defer sqlDB.Close()

	// Initialize handlers
	httpRouter := router.NewMuxRouter()

	pingHandler := handler.NewPingHandler()
	httpRouter.GET("/ping", pingHandler.Ping)

	// Initialize product service
	productRepo := exRepo.NewProductRepository(cfg.ExternalConnection.ProductService.Host)
	orderRepository := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepository, productRepo)
	orderHandler := handler.NewOrderHandler(orderService)
	httpRouter.POST("/order", orderHandler.CreateOrderHandler)
	httpRouter.GET("/order/{orderID}", orderHandler.GetOrderHandler)

	httpRouter.SERVE(cfg.AppPort)
}

func initLogging(logFolder string) {
	logging.InitLogger()
	currentDate := time.Now().Format("2006-01-02")

	// Specify the log file with the current date
	logFilePath := fmt.Sprintf("%s/app_%s.log", logFolder, currentDate)

	// Create the log file if it doesn't exist
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}

	// Set the logrus output to the log file
	logging.Log.SetOutput(logFile)

	go func() {
		for {
			time.Sleep(time.Hour) // Adjust the sleep duration as needed
			newDate := time.Now().Format("2006-01-02")
			if newDate != currentDate {
				currentDate = newDate
				logFilePath = fmt.Sprintf("%s/app_%s.log", logFolder, currentDate)
				newLogFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					logging.Log.Fatal("Error creating log file:", err)
				}
				logFile = newLogFile
				logging.Log.SetOutput(logFile)
			}
		}
	}()
}
