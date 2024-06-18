package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/lyteabovenyte/microservices-main/order/internal/application/core/domain"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type OrderDatabaseTestSuite struct {
	suite.Suite          // enables the test suite
	DataSourceUrl string // data source url for each test
}

func (o *OrderDatabaseTestSuite) SetupSuite() {
	ctx := context.Background()
	port := "3306/tcp"
	dbUrl := func(host string, port nat.Port) string {
		return fmt.Sprintf("root:s3cr3t@tcp(%s:%s)/orders?charset=utf8mb4&parseTime=True&loc=Local",
			host,
			port.Port())
	}
	req := testcontainers.ContainerRequest{
		Image:        "docker.io/mysql:8.0.30",
		ExposedPorts: []string{port},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "s3cr3t",
			"MYSQL_DATABASE":      "orders",
		},
		WaitingFor: wait.ForSQL(nat.Port(port), "mysql", dbUrl),
	}
	mysqlContainer, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		log.Fatal("Failed to start mysql.", err)
	}
	endpoint, _ := mysqlContainer.Endpoint(ctx, "")
	o.DataSourceUrl = fmt.Sprintf("root:s3cr3t@tcp(%s)/orders?charset=utf8mb4&parseTime=True&loc=Local", endpoint)

}

func (o *OrderDatabaseTestSuite) Test_Should_Save_Orders() {
	adapter, err := NewAdapter(o.DataSourceUrl)
	o.Nil(err)
	saveErr := adapter.Save(&domain.Order{})
	o.Nil(saveErr)
}

func (o *OrderDatabaseTestSuite) Test_Should_Get_Orders() {
	adapter, _ := NewAdapter(o.DataSourceUrl)
	order := domain.NewOrder(2, []domain.OrderItem{
		{
			ProductCode: "CAM",
			UnitPrice:   1.5,
			Quantity:    2,
		},
	})
	adapter.Save(&order)
	ord, err := adapter.Get(fmt.Sprint(order.ID))
	o.Nil(err)
	o.Equal(int64(2), ord.CustomerID)
}

// get all of testingSuite (in our case OrderDatabaseTestSuite) and run all the tests attached to it.
func TestOrderDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderDatabaseTestSuite))
}
