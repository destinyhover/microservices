package api

import (
	"context"
	"errors"
	"testing"

	"github.com/destinyhover/microservices/order/internal/application/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func Test_Return_DB_Fail_Error(t *testing.T) {
	payment := new(mockPayment)
	db := new(mockDb)
	db.On("Save", mock.Anything, mock.Anything).Return(errors.New("db error"))
	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)

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

}
func Test_Return_Payment_Fail_Error(t *testing.T) {
	payment := new(mockPayment)
	db := new(mockDb)
	db.On("Save", mock.Anything, mock.Anything).Return(nil)
	payment.On("Charge", mock.Anything, mock.Anything).Return(errors.New("insufficient balance"))

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
	st, _ := status.FromError(err)
	assert.Equal(t, st.Message(), "order creation failed")
	assert.Equal(t, st.Code(), codes.InvalidArgument)
	det := st.Details()[0].(*errdetails.BadRequest).FieldViolations[0].Description
	assert.Equal(t, det, "insufficient balance")
}
