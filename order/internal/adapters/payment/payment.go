package payment

import (
	"context"
	"github.com/huseyinbabal/microservices-proto/golang/payment"
	"github.com/lyteabovenyte/microservices-main/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	payment payment.PaymentClient
}

func NewAdapter(PaymentServiceURL string) (*Adapter, error) {
	var opts []grpc.DialOption // data model for connection configuration.
	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(PaymentServiceURL, opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := payment.NewPaymentClient(conn)
	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(o *domain.Order) error {
	_, err := a.payment.Create(context.Background(),
		&payment.CreatePaymentRequest{
			OrderId:    o.ID,
			UserId:     o.CustomerID,
			TotalPrice: o.TotalPrice()})
	return err
}
