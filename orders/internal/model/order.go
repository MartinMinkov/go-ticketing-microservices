package model

import (
	"context"
	"errors"
	"time"

	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Status string

const (
	OrderCreated         Status = "created"
	OrderCancelled       Status = "cancelled"
	OrderAwaitingPayment Status = "awaiting:payment"
	OrderComplete        Status = "complete"
)

const DefaultDuration = time.Minute * 2

type Order struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UserId    *string            `bson:"user_id" json:"user_id" validate:"required"`
	TicketId  *string            `bson:"ticket_id" json:"ticket_id" validate:"required"`
	Status    *string            `bson:"status" json:"status"`
	ExpiresAt *time.Time         `bson:"expires_at" json:"expires_at"`
	Version   *int64             `bson:"version" json:"version"`
}

func NewOrder(userId string, ticketId string) *Order {
	defaultVersion := int64(0)
	defaultStatus := string(OrderCreated)
	defaultExpiresAt := time.Now().Add(DefaultDuration)
	return &Order{
		ID:        primitive.NewObjectID(),
		UserId:    &userId,
		TicketId:  &ticketId,
		Status:    &defaultStatus,
		ExpiresAt: &defaultExpiresAt,
		Version:   &defaultVersion,
	}
}

func (o *Order) incrementVersion() int64 {
	return *o.Version + int64(1)
}

func GetAllOrders(db *database.Database, userId string) ([]*Order, error) {
	var orders []*Order
	cursor, err := db.OrderCollection.Find(db.Ctx, bson.D{{Key: "user_id", Value: userId}})
	if err != nil {
		return nil, errors.New("failed to get orders: " + err.Error())
	}
	defer cursor.Close(db.Ctx)
	for cursor.Next(db.Ctx) {
		var order Order
		err := cursor.Decode(&order)
		if err != nil {
			return nil, errors.New("failed to decode order: " + err.Error())
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

func GetSingleOrder(db *database.Database, id string, userId string) (*Order, error) {
	var order Order
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("failed to order ticket ID to ObjectID: " + err.Error())
	}
	err = db.OrderCollection.FindOne(db.Ctx, bson.D{{Key: "_id", Value: objectId}, {Key: "user_id", Value: userId}}).Decode(&order)
	if err != nil {
		return nil, errors.New("failed to get order: " + err.Error())
	}
	return &order, nil
}

func (o *Order) Save(db *database.Database) error {
	_, err := db.OrderCollection.InsertOne(db.Ctx, o)
	if err != nil {
		return errors.New("failed to save order: " + err.Error())
	}
	return nil
}

func (o *Order) Update(db *database.Database) error {
	filter := bson.M{"_id": o.ID, "version": o.Version}
	update := bson.M{
		"$set": bson.M{
			"status":     o.Status,
			"expires_at": o.ExpiresAt,
			"version":    o.incrementVersion(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedOrder Order
	err := db.OrderCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedOrder)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("order not found")
		}
		return errors.New("failed to update order: " + err.Error())
	}

	*o = updatedOrder
	return nil
}

func CancelOrder(db *database.Database, id string) (*Order, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("failed to order ticket ID to ObjectID: " + err.Error())
	}
	filter := bson.M{"_id": objectId}
	update := bson.M{
		"$set": bson.M{
			"status": OrderCancelled,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedOrder Order
	err = db.OrderCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedOrder)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("order not found")
		}
		return nil, errors.New("failed to update order: " + err.Error())
	}

	return &updatedOrder, nil
}

func GetOrderByTicketId(db *database.Database, ticketId string) (*Order, error) {
	var order Order
	err := db.OrderCollection.FindOne(db.Ctx, bson.D{{Key: "ticket_id", Value: ticketId}}).Decode(&order)
	if err != nil {
		return nil, errors.New("failed to get order: " + err.Error())
	}
	return &order, nil
}

func IsTicketReserved(db *database.Database, ticketId string) (bool, error) {
	var order Order
	filter := bson.M{"$in": bson.A{string(OrderCreated), string(OrderAwaitingPayment), string(OrderComplete)}}
	err := db.OrderCollection.FindOne(db.Ctx, bson.D{{Key: "ticket_id", Value: ticketId}, {Key: "status", Value: filter}}).Decode(&order)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, errors.New("failed to get order: " + err.Error())
	}
	return true, nil
}
