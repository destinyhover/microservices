package db

import (
	"context"
	"fmt"

	"github.com/destinyhover/microservices/order/internal/application/core/domain"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrderItem struct {
	gorm.Model
	ProductCode string
	UnitPrice   float32
	Quantity    int32
	OrderID     uint
}
type Order struct {
	gorm.Model
	CustomerId int64
	Status     string
	OrderItems []OrderItem
}
type Adapter struct {
	db *gorm.DB
}

func (a Adapter) Get(ctx context.Context, id int64) (domain.Order, error) {
	var orderEnt Order
	var orderItems []domain.OrderItem
	res := a.db.WithContext(ctx).Preload("OrderItems").First(&orderEnt, id)

	for _, v := range orderEnt.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{ProductCode: v.ProductCode, UnitPrice: v.UnitPrice,
			Quantity: v.Quantity})
	}

	result := domain.Order{ID: int64(orderEnt.ID), CustomerId: orderEnt.CustomerId, Status: orderEnt.Status,
		OrderItems: orderItems, CreatedAt: orderEnt.CreatedAt.Unix()}
	return result, res.Error
}

func (a Adapter) Save(ctx context.Context, order *domain.Order) error {
	var orderItems []OrderItem
	for _, v := range order.OrderItems {
		orderItems = append(orderItems, OrderItem{
			ProductCode: v.ProductCode,
			UnitPrice:   v.UnitPrice,
			Quantity:    v.Quantity,
		})
	}
	orderEnt := Order{
		CustomerId: order.CustomerId,
		Status:     order.Status,
		OrderItems: orderItems,
	}

	res := a.db.WithContext(ctx).Create(&orderEnt)
	if res.Error != nil {
		order.ID = int64(orderEnt.ID)
		order.CreatedAt = orderEnt.CreatedAt.Unix()
		return res.Error
	}
	order.ID = int64(orderEnt.ID)
	order.CreatedAt = orderEnt.CreatedAt.Unix()
	return nil
}

func NewAdapter(sourceUrl string) (*Adapter, error) {
	db, err := gorm.Open(postgres.Open(sourceUrl), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db connection err: %w", err)
	}

	if err = db.Use(otelgorm.NewPlugin(otelgorm.WithDBName("order"))); err != nil {
		return nil, fmt.Errorf("db otel plugin err: %w", err)
	}

	if err := db.AutoMigrate(&Order{}, &OrderItem{}); err != nil {
		return nil, fmt.Errorf("db migration err: %w", err)
	}

	return &Adapter{db: db}, nil
}
