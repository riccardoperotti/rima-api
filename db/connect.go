package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// TODO: use some kind of secret management
const (
	username = "verso"
	password = "verso$1"
	hostname = "localhost"
	dbname   = "verso"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error when opening DB: %s", err)
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors pinging DB: %s", err)
		return nil, err
	}

	log.Printf("Connected to DB %s successfully\n", dbname)

	return db, nil
}
