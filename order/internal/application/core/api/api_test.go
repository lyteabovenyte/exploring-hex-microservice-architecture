package api

import (
	"context"
	"testing"

	"github.com/lyteabovenyte/microservices-main/order/internal/application/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedPayment struct {
	mock.Mock
}

func (p *mockedPayment) Charge(ctx context.Context, order *domain.Order) error {
	args := p.Called(ctx, order)
	return args.Error(0)
}

type mockedDB struct {
	mock.Mock
}

func (d *mockedDB) Save(ctx context.Context, order *domain.Order) error {
	args := d.Called(ctx, order)
	return args.Error(0)
}

func (d *mockedDB) Get(ctx context.Context, id string) (domain.Order, error) {
	args := d.Called(ctx, id)
	return args.Get(0).(domain.Order), args.Error(1)
}

func TestPlaceOrder(t *testing.T) {
	payment := new(mockedPayment)
	db := new(mockedDb)
	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)
	db.On("Save", mock.Anything, mock.Anything).Return(nil)

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(domain.Order{
		CustomerID: 123,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "camera",
				UnitPrice:   12.3,
				Quantity:    3,
			},
		},
		CreatedAt: 0,
	})
	assert.Nil(t, err)
}


