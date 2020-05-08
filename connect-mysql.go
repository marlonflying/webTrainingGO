package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	connHost       = "localhost"
	connPort       = "8080"
	driverName     = "mysql"
	dataSourceName = "root:datapassword@/mydb"
)

var db *sql.DB
var connError error

func init() {
	db, connError = sql.Open(driverName, dataSourceName)
	if connError != nil {
		log.Fatal("error connecting to database :: ", connError)
	}
}

func getCurrentDb(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT DATABASE() AS db")
	if err != nil {
		log.Print("error executing query :: ", err)
		return
	}
	var db string
	for rows.Next() {
		rows.Scan(&db)
	}
	fmt.Fprintf(w, "Current database is :: %s", db)
}

func main() {
	http.HandleFunc("/", getCurrentDb)
	defer db.Close()
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}
}
