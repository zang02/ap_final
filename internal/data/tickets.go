package data

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Ticket struct {
	ID        int64     `json:"id"`
	UserId    int64     `json:"userId"`
	CreatedAt string    `json:"time"`
	Total     int64     `json:"total"`
	Products  []Product `json:"products"`
}

type Product struct {
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Amount int    `json:"amount"`
}

type TicketModel struct {
	DB *mongo.Database
}
