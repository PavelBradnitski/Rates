package repositories

import (
	"database/sql"
	"fmt"
	"log"
)

func ConnectToDB(dbUser, dbPassword, dbHost, dbPort, dbName string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting database: %v\n", err)
	}
	return db, nil
}
