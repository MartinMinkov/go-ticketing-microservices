package database

import (
	"context"
	"log"
	"time"

	"github.com/MartinMinkov/go-ticketing-microservices/orders/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client           *mongo.Client
	OrderCollection  *mongo.Collection
	TicketCollection *mongo.Collection
	Ctx              context.Context
	CancelFunc       context.CancelFunc
}

func ConnectDB(config *config.Config) *Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var client *mongo.Client
	var err error
	if config.ApplicationConfig.Environment == "development" {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.DatabaseConfig.GetConnectionStringWithUser()+"?authSource=admin"))
	} else {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.DatabaseConfig.GetConnectionString()))
	}

	if err != nil {
		panic("Failed to connect to MongoDB")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	orderCollection := client.Database(config.DatabaseConfig.Database).Collection("orders")
	ticketCollection := client.Database(config.DatabaseConfig.Database).Collection("tickets")
	return &Database{Client: client, OrderCollection: orderCollection, TicketCollection: ticketCollection, Ctx: context.TODO()}
}
