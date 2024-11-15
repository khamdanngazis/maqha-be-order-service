// internal/app/repository/mock/mock_product_repository.go

package mock

import (
	"context"
	"errors"
	"maqhaa/order_service/external/entity"
)

// MockProductRepository is a mock implementation of the ProductRepository interface.
type MockProductRepository struct {
	GetProductByIDFunc func(ctx context.Context, productID uint, token string) (*entity.Product, error)

	// Map to store dynamic responses for different product IDs
	productResponses map[uint]*entity.Product
}

// NewMockProductRepository creates a new instance of MockProductRepository with an empty response map.
func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		productResponses: make(map[uint]*entity.Product),
	}
}

// SetProductResponse sets a dynamic response for a specific product ID.
func (m *MockProductRepository) SetProductResponse(productID uint, product *entity.Product) {
	m.productResponses[productID] = product
}

// GetProductByID is the mock implementation for the GetProductByID method.
func (m *MockProductRepository) GetProductByID(ctx context.Context, productID uint, token string) (*entity.Product, error) {
	if m.GetProductByIDFunc != nil {
		return m.GetProductByIDFunc(ctx, productID, token)
	}

	// Check if there is a dynamic response for the given product ID
	if product, ok := m.productResponses[productID]; ok {
		return product, nil
	}

	return nil, errors.New("GetProductByIDFunc not implemented in the mock")
}
