package e2e

import (
	"context"
	"log"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/huseyinbabal/microservices-proto/golang/order"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CreateOrderTestSuite struct {
	suite.Suite
	compose *tc.LocalDockerCompose
}

func (c *CreateOrderTestSuite) SetupSuite() {
	ComposeFilePath := []string{"resources/docker-compose.yml"}
	identifier := strings.ToLower(uuid.New().String()) // randomized docker-compose file name
	compose := tc.NewLocalDockerCompose(ComposeFilePath, identifier)

	c.compose = compose

	execErr := compose.WithCommand([]string{"up", "-d"}).Invoke()
	if execErr.Error != nil {
		log.Fatalf("cannot run compose stack %v", execErr.Error)
	}
	// err := execErr.Error()
	// if err != nil {
	// 	log.Fatalf("cannot run compose stack %v", err)
	// }
}

func (c *CreateOrderTestSuite) Test_Should_Create_Order() {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("failed to connect to order service. err: %v", err)
	}
	orderClient := order.NewOrderClient(conn)
	// testing order's Create Mthod.
	CreateOrderResponse, errCreate := orderClient.Create(context.Background(),
		&order.CreateOrderRequest{
			UserId: 23,
			OrderItems: []*order.OrderItem{
				{
					ProductCode: "CAM Canon",
					UnitPrice:   1.5,
					Quantity:    3,
				},
			},
		})
	c.Nil(errCreate)
	// testing order's Get method.
	getOrderResponse, errGet := orderClient.Get(context.Background(),
		&order.GetOrderRequest{OrderId: CreateOrderResponse.OrderId})
	c.Nil(errGet)
	orderItem := getOrderResponse.OrderItems[0]
	c.Equal(float32(1.5), orderItem.ProductCode)
	c.Equal(int32(3), orderItem.Quantity)
	c.Equal("CAM Canon", orderItem.ProductCode)
}

func (c *CreateOrderTestSuite) Tear_Down_Suite() {
	execErr := c.compose.WithCommand([]string{"down"}).Invoke()
	if execErr.Error != nil {
		log.Fatalf("cannot shut the compose stack down. Error: %v", execErr.Error)
	}
	// err := execErr.Error()
	// if err != nil {
	// 	log.Fatalf("cannot shut the compose stack down. Error: %v", err)
	// }
}

// runner
func TestCreateOrderTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderTestSuite))
}
