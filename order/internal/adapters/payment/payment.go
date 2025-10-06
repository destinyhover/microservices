package payment

import (
	"context"
	"log/slog"
	"time"

	pb "github.com/destinyhover/microservices-proto/golang/payment"
	"github.com/destinyhover/microservices/order/internal/application/core/domain"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	client pb.PaymentClient
	conn   *grpc.ClientConn
}

func NewAdapter(paymentServiceURL string) (*Adapter, error) {
	conn, err := grpc.NewClient(
		paymentServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		return nil, err
	}

	return &Adapter{client: pb.NewPaymentClient(conn), conn: conn}, nil
}

func (a *Adapter) Close() error {
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

func (a *Adapter) Charge(ctx context.Context, order *domain.Order) error {
	cctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := a.client.Create(cctx, &pb.PaymentCreateRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})
	if err != nil {
		slog.ErrorContext(cctx, "Charge()", "err", err)
		return err
	}
	return nil
}
