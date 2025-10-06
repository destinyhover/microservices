package api

import (
	"errors"
	"testing"

	"github.com/destinyhover/microservices/order/internal/application/core/domain"
	mocks "github.com/destinyhover/microservices/order/mocks/internal_/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPlaceOrder_paymentFail(t *testing.T) {
	payment := mocks.NewPaymentPort(t)
	db := mocks.NewDBPort(t)
	db.On("Save", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)

	payment.On("Charge", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(errors.New("insufficient balance"))
	app := NewApplication(db, payment)

	_, err := app.PlaceOrder(t.Context(), domain.Order{CustomerID: 1,
		OrderItems: []domain.OrderItem{{ProductCode: "Camera", Quantity: 1, UnitPrice: 30}}})

	require.Error(t, err)
}
func Test_Return_DB_Error(t *testing.T) {
	payment := mocks.NewPaymentPort(t)
	db := mocks.NewDBPort(t)
	db.On("Save", mock.Anything, mock.Anything).Return(errors.New("db error"))
	//payment.On("Charge", mock.Anything, mock.Anything).Return(nil)

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

	assert.EqualError(t, err, "db error")
	payment.AssertNotCalled(t, "Charge", mock.Anything, mock.Anything)
}
