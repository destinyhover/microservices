package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/destinyhover/microservices-proto/golang/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CreateOrderTestSuite struct {
	suite.Suite
	stack compose.ComposeStack
}

func (c *CreateOrderTestSuite) SetupSuite() {

	var err error
	c.stack, err = compose.NewDockerComposeWith(compose.WithStackFiles("resources/docker-compose.yml"), compose.StackIdentifier("pg_e2e_"+uuid.NewString()))
	c.Require().NoError(err)

	c.stack = c.stack.WaitForService("order", wait.ForListeningPort("8181/tcp"))
	c.Require().NoError(c.stack.Up(context.Background(), compose.Wait(true)))
}

func (c *CreateOrderTestSuite) Test_Should_Create_Order() {

	conn, err := grpc.NewClient("localhost:8181", grpc.WithTransportCredentials(insecure.NewCredentials()))

	c.Require().NoError(err, "failed to connect order service")
	defer conn.Close()

	orderClient := order.NewOrderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	createOrderResponse, errCreate := orderClient.Create(ctx, &order.CreateOrderRequest{
		UserId: 23,
		OrderItems: []*order.OrderItem{
			{
				ProductCode: "CAM123",
				Quantity:    3,
				UnitPrice:   1.23,
			},
		},
	})
	c.Require().NoError(errCreate)

	getOrderResponse, errGet := orderClient.Get(ctx, &order.GetOrderRequest{OrderId: createOrderResponse.OrderId})
	c.Require().NoError(errGet)

	c.Equal(int64(23), getOrderResponse.UserId)
	c.Require().Len(getOrderResponse.OrderItems, 1)

	item := getOrderResponse.OrderItems[0]
	c.Equal(float32(1.23), item.UnitPrice)
	c.Equal(int32(3), item.Quantity)
	c.Equal("CAM123", item.ProductCode)
}

func (c *CreateOrderTestSuite) TearDownSuite() {
	c.Require().NoError(
		c.stack.Down(context.Background(), compose.RemoveVolumes(true)),
	)
}

func TestCreateOrderTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderTestSuite))
}
