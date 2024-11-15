// internal/repository/product_repository.go

package repository

import (
	"context"
	"errors"
	"maqhaa/library/logging"
	"maqhaa/library/middleware"
	"maqhaa/order_service/external/entity"
	pb "maqhaa/order_service/external/model"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// ProductRepository handles database interactions related to products.
type ProductRepository interface {
	GetProductByID(ctx context.Context, productID uint, token string) (*entity.Product, error)
}

// Implement the interface in the ProductRepository struct
type productRepository struct {
	connetionURl string
}

// NewProductRepository creates a new ProductRepository instance.
func NewProductRepository(connetionURl string) ProductRepository {
	return &productRepository{connetionURl: connetionURl}
}

func (r *productRepository) GetProductByID(ctx context.Context, productID uint, token string) (*entity.Product, error) {
	requestID, _ := ctx.Value(middleware.RequestIDKey).(string)
	req := &pb.GetProductRequest{
		ProductId: uint32(productID),
		Token:     token, // Replace with a valid product ID for your test data
	}
	conn, err := grpc.Dial(r.connetionURl, grpc.WithInsecure())
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error GetProductByID  %s", err.Error())
		return nil, err
	}
	defer conn.Close()

	client := pb.NewProductClient(conn)

	resp, err := client.GetProduct(context.Background(), req)
	if err != nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error GetProductByID  %s", err.Error())
		return nil, err
	}
	if resp.Code != 0 {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error GetProductByID  %v", resp)
		return nil, errors.New(resp.Message)
	}

	if resp.Data == nil {
		logging.Log.WithFields(logrus.Fields{"request_id": requestID}).Errorf("Error GetProductByID  %s", errors.New("Data Product Nill"))
		return nil, errors.New("Data Product Nill")
	}

	product := &entity.Product{
		ID:          uint(resp.Data.Id),
		Name:        resp.Data.Name,
		Image:       resp.Data.Image,
		Price:       float64(resp.Data.Price),
		Description: resp.Data.Description,
		IsActive:    resp.Data.IsActive,
		CreatedAt:   resp.Data.CreatedAt,
	}

	return product, nil
}
