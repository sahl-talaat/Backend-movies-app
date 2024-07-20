package config

import (
	"context"
	"fmt"
	"log"
	"myapp/ent"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var Client *ent.Client

func InitDB() error {
	client, err := ent.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}

	err = client.Schema.Create(context.Background())
	if err != nil {
		log.Fatalf("failed creating schema resources: %v:", err)
	}

	Client = client
	return nil
}
