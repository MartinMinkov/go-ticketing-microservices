package model

import (
	"context"
	"errors"

	"github.com/MartinMinkov/go-ticketing-microservices/tickets/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Ticket struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	UserId  *string            `bson:"user_id" json:"user_id" validate:"required"`
	Title   *string            `bson:"title" json:"title" validate:"required"`
	Price   *int64             `bson:"price" json:"price"`
	Version *int64             `bson:"version" json:"version"`
}

func NewTicket(userId string, title string, price int64) *Ticket {
	defaultVersion := int64(0)
	return &Ticket{
		ID:      primitive.NewObjectID(),
		UserId:  &userId,
		Title:   &title,
		Price:   &price,
		Version: &defaultVersion,
	}
}

func (t *Ticket) incrementVersion() int64 {
	return *t.Version + int64(1)
}

func GetAllTickets(db *database.Database, userId string) ([]*Ticket, error) {
	var tickets []*Ticket
	cursor, err := db.TicketCollection.Find(db.Ctx, bson.D{{Key: "user_id", Value: userId}})
	if err != nil {
		return nil, errors.New("failed to get tickets: " + err.Error())
	}
	defer cursor.Close(db.Ctx)
	for cursor.Next(db.Ctx) {
		var ticket Ticket
		err := cursor.Decode(&ticket)
		if err != nil {
			return nil, errors.New("failed to decode ticket: " + err.Error())
		}
		tickets = append(tickets, &ticket)
	}
	return tickets, nil
}

func GetSingleTicket(db *database.Database, id string) (*Ticket, error) {
	var ticket Ticket
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("failed to convert ticket ID to ObjectID: " + err.Error())
	}
	err = db.TicketCollection.FindOne(db.Ctx, bson.D{{Key: "_id", Value: objectId}}).Decode(&ticket)
	if err != nil {
		return nil, errors.New("failed to get ticket: " + err.Error())
	}
	return &ticket, nil
}

func GetSingleTicketByTitle(db *database.Database, title string) (*Ticket, error) {
	var ticket Ticket
	err := db.TicketCollection.FindOne(db.Ctx, bson.D{{Key: "title", Value: title}}).Decode(&ticket)
	if err != nil {
		return nil, errors.New("failed to get ticket: " + err.Error())
	}
	return &ticket, nil
}

func (t *Ticket) Save(db *database.Database) error {
	_, err := db.TicketCollection.InsertOne(db.Ctx, t)
	if err != nil {
		return errors.New("failed to save ticket: " + err.Error())
	}
	return nil
}

func (t *Ticket) Update(db *database.Database) error {
	filter := bson.M{"_id": t.ID, "version": t.Version}
	update := bson.M{
		"$set": bson.M{
			"title":   t.Title,
			"price":   t.Price,
			"version": t.incrementVersion(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedTicket Ticket
	err := db.TicketCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedTicket)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("ticket not found")
		}
		return errors.New("failed to update ticket: " + err.Error())
	}

	*t = updatedTicket
	return nil
}

func (t *Ticket) Delete(db *database.Database) error {
	_, err := db.TicketCollection.DeleteOne(db.Ctx, bson.D{{Key: "_id", Value: t.ID}})
	if err != nil {
		return errors.New("failed to delete ticket: " + err.Error())
	}
	return nil
}
