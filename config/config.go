package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost    string
	JWTSecret string
}

var db *pgx.Conn

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return &Config{
		DBHost:    os.Getenv("POSTGRES_CONNECTION_STRING"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}, nil
}

func InitDB() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	db = conn
}

func GetDB() *pgx.Conn {
	return db
}
