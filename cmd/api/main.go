package main

import (
	"app/internal/data"
	"app/internal/jsonlog"
	"context"
	"os"
	"text/template"
	"time"
)

type application struct {
	config        config
	models        data.Models
	logger        *jsonlog.Logger
	templateCache map[string]*template.Template
}

type config struct{}

func main() {
	godotenv.Load()
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DATABASE_URI")))
	if err != nil {
		if err != nil {
			return nil, err
		}
	}
	return client
}
