package fakes

import (
	"context"

	"github.com/destinyhover/microservices/order/internal/application/core/domain"
	"github.com/destinyhover/microservices/order/internal/ports"
)

// NoopPayment — заглушка, которая всегда "успешно" проводит платеж.
type NoopPayment struct{}

// Проверка на этапе компиляции: NoopPayment реализует PaymentPort
var _ ports.PaymentPort = (*NoopPayment)(nil)

func (NoopPayment) Charge(ctx context.Context, order *domain.Order) error {
	// Ничего не делаем, просто возвращаем nil
	return nil
}
