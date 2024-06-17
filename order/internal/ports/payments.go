package ports

import "github.com/lyteabovenyte/microservices-main/order/internal/application/core/domain"

// PaymentPort allows the order service to call the payment service.
type PaymentPort interface {
	Charge(*domain.Order) error
}
