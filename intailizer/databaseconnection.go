package intailizer

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func DatabaseConnection() {
	dsn := os.Getenv("DB_URL")

	if dsn == "" {
		log.Fatal("error in loading env in database connection")
	}

	//open connection

	conn, error := sql.Open("postgres", dsn)

	if error != nil {
		log.Print("error in openign datavase connection")
	}

	//check ping
	err := conn.Ping()

	if err != nil {
		log.Print("failed to conected db ")
	}
	fmt.Print("sucessfully connected to db")

	DB = conn

}
