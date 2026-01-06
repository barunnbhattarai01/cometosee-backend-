package intailizer

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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

	//connection pooling
	conn.SetMaxOpenConns(25)                 //maximunnm total conn
	conn.SetMaxIdleConns(10)                 //maximum idle conn
	conn.SetConnMaxLifetime(5 * time.Minute) //timeee to reusee

	//check ping
	err := conn.Ping()

	if err != nil {
		log.Print("failed to conected db ")
	}
	fmt.Print("sucessfully connected to db")

	DB = conn

}
