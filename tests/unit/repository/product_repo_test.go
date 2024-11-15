package repository_test

import (
	"maqhaa/order_service/internal/app/repository/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

var productRepo *mock.MockProductRepository

func TestProductRepo_GetProduct(t *testing.T) {
	client := SampleClient()

	categories := SampleCategories(client.ID)
	producRepo.SetProductResponse(categories[0].Products[0].ID, &categories[0].Products[0])
	producRepo.SetProductResponse(categories[0].Products[1].ID, &categories[0].Products[1])
	producRepo.SetProductResponse(categories[1].Products[0].ID, &categories[1].Products[0])
	producRepo.SetProductResponse(categories[2].Products[0].ID, &categories[2].Products[0])

	product, err := producRepo.GetProductByID(ctx, categories[0].Products[0].ID, client.Token)

	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, categories[0].Products[0].ID, product.ID)

	product1, err := producRepo.GetProductByID(ctx, categories[0].Products[1].ID, client.Token)

	assert.Nil(t, err)
	assert.NotNil(t, product1)
	assert.Equal(t, categories[0].Products[1].ID, product1.ID)

	product3, err := producRepo.GetProductByID(ctx, categories[1].Products[0].ID, client.Token)

	assert.Nil(t, err)
	assert.NotNil(t, product3)
	assert.Equal(t, categories[1].Products[0].ID, product3.ID)

	product4, err := producRepo.GetProductByID(ctx, categories[2].Products[0].ID, client.Token)

	assert.Nil(t, err)
	assert.NotNil(t, product4)
	assert.Equal(t, categories[2].Products[0].ID, product4.ID)
}
