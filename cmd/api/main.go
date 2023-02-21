package main

import (
	"app/internal/data"
	"app/internal/woodlog"
	"context"
	"flag"
	"html/template"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type application struct {
	config        config
	models        data.Models
	logger        *woodlog.Logger
	templateCache map[string]*template.Template

	wg sync.WaitGroup
}

type config struct {
	port string
	db   struct {
		dns string
	}
}

func main() {
	godotenv.Load()
	var config config
	flag.StringVar(&config.port, "port", os.Getenv("PORT"), "port")
	flag.StringVar(&config.db.dns, "uri", os.Getenv("MONGOURI"), "mongo uri")

	logger := *woodlog.New(os.Stdout, 0)

	templateCache, err := data.NewTemplateCache("./ui/html/")
	if err != nil {
		logger.PrintFatal(err, "failed to create template cache")
	}

	db, err := openDB(config.db.dns)
	if err != nil {
		logger.PrintFatal(err, "failed to connect to mongo")
	}
	defer db.Disconnect(context.TODO())

	app := application{
		templateCache: templateCache,
		config:        config,
		logger:        &logger,
		models:        data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, "failed to start server")
	}
}

func openDB(dsn string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}
	return client, nil
}
