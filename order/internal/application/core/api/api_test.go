package api

import (
	"context"
	"testing"

	"github.com/destinyhover/microservices/order/internal/application/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPayment struct {
	mock.Mock
}

func (p *mockPayment) Charge(ctx context.Context, order *domain.Order) error {
	args := p.Called(ctx, order)
	return args.Error(0)
}

type mockDb struct {
	mock.Mock
}

func (db *mockDb) Save(ctx context.Context, order *domain.Order) error {
	args := db.Called(ctx, order)
	return args.Error(0)
}

func (d *mockDb) Get(ctx context.Context, id int64) (domain.Order, error) {
	args := d.Called(ctx, id)
	return args.Get(0).(domain.Order), args.Error(0)
}

func TestPlaceOrder(t *testing.T) {
	payment := new(mockPayment)
	db := new(mockDb)
	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)
	db.On("Save", mock.Anything, mock.Anything).Return(nil)

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(t.Context(), domain.Order{
		CustomerID: 1234,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "Camera",
				UnitPrice:   12.3,
				Quantity:    3,
			},
		},
		CreatedAt: 0,
	})
	assert.NoError(t, err)
}
