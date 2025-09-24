package grpc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/destinyhover/microservices-proto/golang/payment"
	"github.com/destinyhover/microservices/payment/internal/application/core/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a Adapter) Create(ctx context.Context, request *payment.PaymentCreateRequest) (*payment.PaymentCreateResponse, error) {
	slog.InfoContext(ctx, "creating payment...")
	newPayment := domain.NewPayment(request.UserId, request.OrderId, int32(request.TotalPrice))

	result, err := a.api.Charge(ctx, newPayment)
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to charge %v", err)).Err()
	}
	return &payment.PaymentCreateResponse{PaymentId: result.ID}, nil
}
