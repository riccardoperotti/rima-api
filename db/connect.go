package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// These must come from config:
// TODO: use some kind of secret management
const (
	username = "lex"
	password = "$l3x1c4"
	hostname = "localhost"
	dbname   = "lexico"
)

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

func Connect() (*gorm.DB, error) {

	dsn := dsn()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Printf("Error when opening DB: %s", err)
		return nil, err
	}

	log.Printf("Connected to DB %s successfully\n", dbname)

	return db, nil
}
