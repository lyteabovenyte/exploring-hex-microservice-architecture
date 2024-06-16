package api

import ( 
	"github.com/lyteabovenyte/microservices-main/order/internal/application/core/domain"
	"github.com/lyteabovenyte/microservices-main/order/internal/ports"
)

// the API depends on DBPort
type Application struct {
	db ports.DBPort
}

// DBPort is passed during the app initialization
func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

// order is saved through the DBPort
func (a *Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	err := a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}
	return order, nil
}