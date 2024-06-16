// ports are generally interfaces that contains information about
// interaction between an actor and an application.
package ports

import (
	"github.com/lyteabovenyte/microservices-main/order/internal/application/core/domain"
)

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(*domain.Order) error
}