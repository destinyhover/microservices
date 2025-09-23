package db

import (
	"context"
	"fmt"
	"time"

	"github.com/destinyhover/microservices/payment/internal/application/core/domain"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	CustomerID int64
	Status     string
	OrderID    int64
	TotalPrice int32
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(strDb string) (*Adapter, error) {
	db, err := gorm.Open(postgres.Open(strDb), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connection db err: %v", err)
	}
	err = db.Use(otelgorm.NewPlugin(otelgorm.WithDBName("payment")))
	if err != nil {
		return nil, fmt.Errorf("db otel plugin err: %v", err)
	}

	if err = db.AutoMigrate(&Payment{}); err != nil {
		return nil, fmt.Errorf("db autoMidgrate err: %v", err)
	}
	return &Adapter{db: db}, nil
}

func (a Adapter) Get(ctx context.Context, id string) (domain.Payment, error) {
	var paymentEntity domain.Payment
	res := a.db.WithContext(ctx).First(&paymentEntity, id)
	payment := domain.Payment{
		ID:         int64(paymentEntity.ID),
		CustomerID: paymentEntity.CustomerID,
		Status:     paymentEntity.Status,
		OrderID:    paymentEntity.OrderID,
		TotalPrice: paymentEntity.TotalPrice,
		CreatedAt:  time.Now().Unix(),
	}
	return payment, res.Error
}

func (a Adapter) Save(ctx context.Context, payment *domain.Payment) error {
	orderModel := Payment{
		CustomerID: payment.CustomerID,
		Status:     payment.Status,
		OrderID:    payment.OrderID,
		TotalPrice: payment.TotalPrice,
	}
	res := a.db.WithContext(ctx).Create(&orderModel)
	if res.Error != nil {
		payment.ID = int64(orderModel.ID)
		payment.CreatedAt = orderModel.CreatedAt.Unix()
		return res.Error
	}
	payment.ID = int64(orderModel.ID)
	payment.CreatedAt = orderModel.CreatedAt.Unix()
	return nil
}
