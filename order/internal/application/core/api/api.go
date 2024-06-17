package api

import (
	"strings"

	"github.com/lyteabovenyte/microservices-main/order/internal/application/core/domain"
	"github.com/lyteabovenyte/microservices-main/order/internal/ports"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// the API depends on DBPort
type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort // the order service depends on the payment service to charge the customers.
}

// DBPort is passed during the app initialization
func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

// order is saved through the DBPort
func (a *Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	err := a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}
	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		// Error with details
		st := status.Convert(paymentErr)
		var allErrors []string
		for _, details := range st.Details() {
			switch t := details.(type) {
			case *errdetails.BadRequest:
				for _, violation := range t.GetFieldViolations() {
					allErrors = append(allErrors, violation.Description)
				}
			}
		}
		fieldErr := &errdetails.BadRequest_FieldViolation{
			Field:       "payment",
			Description: strings.Join(allErrors, "\n"),
		}
		badReq := &errdetails.BadRequest{}
		badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)
		orderStatus := status.New(codes.InvalidArgument, "order creation failed")
		statusWithDeatils, _ := orderStatus.WithDetails(badReq)
		return domain.Order{}, statusWithDeatils.Err()
	}
	return order, nil
}
