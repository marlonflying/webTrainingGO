package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	connHost   = "localhost"
	connPort   = "8080"
	mongoDbURL = "127.0.0.1"
)

var session *mgo.Session
var connError error

type Employee struct {
	Id   int    `json:"uid"`
	Name string `json:"name"`
}

func init() {
	session, connError = mgo.Dial(mongoDbURL)
	if connError != nil {
		log.Fatal("error connecting to database :: ", connError)
	}
	session.SetMode(mgo.Monotonic, true)
}

func getDbNames(w http.ResponseWriter, r *http.Request) {
	db, err := session.DatabaseNames()
	if err != nil {
		log.Print("error getting database names :: ", err)
		return
	}
	fmt.Fprintf(w, "Databases names are :: %s", strings.Join(db, ", "))
}

func createDocument(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	name, nameOk := vals["name"]
	id, idOk := vals["id"]
	if nameOk && idOk {
		employeeId, err := strconv.Atoi(id[0])
		if err != nil {
			log.Print("error converting string to int for id :: ", err)
			return
		}
		log.Print("going to insert document in database for name :: ", name[0])
		collection := session.DB("mydb").C("employee")
		err = collection.Insert(&Employee{employeeId, name[0]})
		if err != nil {
			log.Print("error inserting document in database :: ", err)
			return
		}
		fmt.Fprintf(w, "Last created document id is :: %s", id[0])
	} else {
		fmt.Fprintf(w, "Error creating document in database for name :: %s", name[0])
	}
}

func readDocuments(w http.ResponseWriter, r *http.Request) {
	log.Print("reading documents from database")
	var employees []Employee
	collection := session.DB("mydb").C("employee")
	err := collection.Find(bson.M{}).All(&employees)
	if err != nil {
		log.Print("error reading documents from database :: ", err)
		return
	}
	json.NewEncoder(w).Encode(employees)
}

func updateDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	vals := r.URL.Query()
	name, ok := vals["name"]
	if ok {
		employeeId, err := strconv.Atoi(id)
		if err != nil {
			log.Print("error converting string id to int :: ", err)
			return
		}
		log.Print("going to update document in database for id :: ", id)
		collection := session.DB("mydb").C("employee")
		var changeInfo *mgo.ChangeInfo
		changeInfo, err = collection.Upsert(bson.M{"id": employeeId}, &Employee{employeeId, name[0]})
		if err != nil {
			log.Print("error updating record in database :: ", err)
			return
		}
		fmt.Fprintf(w, "Number of documents updated in database are :: %d", changeInfo.Updated)
	} else {
		fmt.Fprintf(w, "Error updating document in database for id :: %s", id)
	}
}

func deleteDocument(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	name, ok := vals["name"]
	if ok {
		log.Print("going to delete document in database for name :: ", name[0])
		collection := session.DB("mydb").C("employee")
		removeErr := collection.Remove(bson.M{"name": name[0]})
		if removeErr != nil {
			log.Print("error removing document form database :: ", removeErr)
			return
		}
		fmt.Fprintf(w, "Document with name %s is deleted from database", name[0])
	} else {
		fmt.Fprintf(w, "error deleting document in database for name %s", name[0])
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/databases", getDbNames).Methods("GET")
	router.HandleFunc("/employee/create", createDocument).Methods("POST")
	router.HandleFunc("/employee/delete", deleteDocument).Methods("DELETE")
	router.HandleFunc("/employees", readDocuments).Methods("GET")
	router.HandleFunc("/employee/update/{id}", updateDocument).Methods("PUT")
	defer session.Close()
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}
}
