package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type User struct {
	Name string
}

type PostgresConfig struct {
	Host string
	Port string
	User string
	Password string
	Database string
	SSLMode string
}

func (cfg PostgresConfig) ToString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode
	)
}

func main() {
	cfg := PostgresConfig{
		Host: "localhost",
		Port: "5435",
		User: "ricardo",
		Password: "ricardo",
		Database: "lens",
		SSLMode: "disable",
	}

	db, err := sql.Open("pgx", cfg.ToString())

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("connected")
}
