//go:build integration

package db

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/destinyhover/microservices/order/internal/application/core/domain"
	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type OrderDBTestSuit struct {
	suite.Suite
	DataSourceURL string
	container     testcontainers.Container
}

func TestOrderDatabaseSuite(t *testing.T) {
	suite.Run(t, new(OrderDBTestSuit))
}
func (o *OrderDBTestSuit) SetupSuite() {
	port := "5432/tcp"
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("postgres://user:pass@%s:%s/postgres?sslmode=disable", host, port.Port())
	}
	container, err := testcontainers.GenericContainer(context.TODO(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16",
			ExposedPorts: []string{port},
			Env: map[string]string{
				"POSTGRES_USER":     "user",
				"POSTGRES_PASSWORD": "pass",
				"POSTGRES_DB":       "orders",
			},
			WaitingFor: wait.ForSQL(nat.Port(port), "pgx", dbURL).WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		slog.Error("SetupSuite: start container", "err", err)
		o.FailNow("failed to start postgres container")
		return
	}
	o.container = container
	endpoint, _ := container.Endpoint(context.Background(), "")
	o.DataSourceURL = fmt.Sprintf("postgres://user:pass@%s/orders?sslmode=disable", endpoint)
}
func (o *OrderDBTestSuit) Test_Sould_Save_Order() {
	adapter, err := NewAdapter(o.DataSourceURL)
	o.Nil(err)
	saveErr := adapter.Save(context.Background(), &domain.Order{})
	o.Nil(saveErr)
}
