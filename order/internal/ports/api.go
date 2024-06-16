// ports are generally interfaces that contains information about
// interaction between an actor and an application.
package ports

import (
	"github.com/lyteabovenyte/microservices-main/order/internal/application/core/domain"
)

// APIPort is used for core application functionalities
type APIPort interface {
	PlaceOrder(order domain.Order) (domain.Order, error)
}
