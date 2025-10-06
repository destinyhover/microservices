package api

import (
	"errors"
	"testing"

	"github.com/destinyhover/microservices/order/internal/application/core/domain"
	mocks "github.com/destinyhover/microservices/order/mocks/internal_/ports"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPlaceOrder_paymentFail(t *testing.T) {
	payment := mocks.NewPaymentPort(t)
	db := mocks.NewDBPort(t)
	payment.On("Charge", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(errors.New("insufficient balance"))
	app := NewApplication(db, payment)

	_, err := app.PlaceOrder(t.Context(), domain.Order{CustomerID: 1,
		OrderItems: []domain.OrderItem{{ProductCode: "Camera", Quantity: 1, UnitPrice: 30}}})

	require.Error(t, err)
	db.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}
