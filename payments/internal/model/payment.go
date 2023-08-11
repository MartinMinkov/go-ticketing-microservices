package model

import (
	"github.com/MartinMinkov/go-ticketing-microservices/payments/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	OrderId  string             `json:"orderId" bson:"order_id"`
	StripeId string             `json:"stripeId" bson:"stripe_id"`
}

func NewPayment(orderId string, stripeId string) *Payment {
	return &Payment{
		ID:       primitive.NewObjectID(),
		OrderId:  orderId,
		StripeId: stripeId,
	}
}

func (p *Payment) Save(db *database.Database) error {
	_, err := db.PaymentCollection.InsertOne(db.Ctx, p)
	if err != nil {
		return err
	}
	return nil
}
