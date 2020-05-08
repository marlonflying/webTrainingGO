package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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

type Employee struct {
	Id   int    `json:"uid"`
	Name string `json:"name"`
}

func createRecord(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	name, ok := vals["name"]
	if ok {
		log.Print("inserting record in database for name : ", name[0])
		stmt, err := db.Prepare("INSERT employee SET name=?")
		if err != nil {
			log.Print("error preparing query :: ", err)
			return
		}
		result, err := stmt.Exec(name[0])
		if err != nil {
			log.Print("error executing query :: ", err)
			return
		}
		id, err := result.LastInsertId()
		fmt.Fprintf(w, "Last inserted Record Id is :: %s", strconv.FormatInt(id, 10))
	} else {
		fmt.Fprintf(w, "Error creating record in database for name :: %s", name[0])
	}
}

func readRecord(w http.ResponseWriter, r *http.Request) {
	log.Print("reading records from database")
	rows, err := db.Query("SELECT * FROM employee")
	if err != nil {
		log.Print("error in select query :: ", err)
		return
	}
	employees := []Employee{}
	for rows.Next() {
		var uid int
		var name string
		err = rows.Scan(&uid, &name)
		employee := Employee{Id: uid, Name: name}
		employees = append(employees, employee)
	}
	json.NewEncoder(w).Encode(employees)
}

func updateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	vals := r.URL.Query()
	name, ok := vals["name"]
	if ok {
		log.Print("updating record in database for id :: ", id)
		stmt, err := db.Prepare("UPDATE employee SET name=? where uid=?")
		if err != nil {
			log.Print("error preparing query :: ", err)
			return
		}
		result, err := stmt.Exec(name[0], id)
		if err != nil {
			log.Print("error executing query :: ", err)
			return
		}
		rowsAffected, err := result.RowsAffected()
		fmt.Fprintf(w, "Number of rows updated :: %d", rowsAffected)
	} else {
		fmt.Fprintf(w, "Error updating record in database for id :: %s", id)
	}
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	name, ok := vals["name"]
	if ok {
		log.Print("deleting record in database for name :: ", name[0])
		stmt, err := db.Prepare("DELETE from employee where name=?")
		if err != nil {
			log.Print("error preparing query :: ", err)
			return
		}
		result, err := stmt.Exec(name[0])
		if err != nil {
			log.Print("error executing query :: ", err)
			return
		}
		rowsAffected, err := result.RowsAffected()
		fmt.Fprintf(w, "Number of rows deleted in database are :: %d", rowsAffected)
	} else {
		fmt.Fprintf(w, "Error deleting record in database for name %s", name[0])
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/employee/create", createRecord).Methods("POST")
	router.HandleFunc("/employees", readRecord).Methods("GET")
	router.HandleFunc("/employee/update/{id}", updateRecord).Methods("PUT")
	router.HandleFunc("/employee/delete", deleteRecord).Methods("DELETE")
	defer db.Close()
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}
}
