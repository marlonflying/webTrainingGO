package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

// Error struct
type NameNotFoundError struct {
	Code int
	Err  error
}

func (nameNotFoundError NameNotFoundError) Error() string {
	return nameNotFoundError.Err.Error()
}

// WrapperHandler
type WrapperHandler func(http.ResponseWriter, *http.Request) error

func (wrapperHandler WrapperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	err := wrapperHandler(w, r)
	if err != nil {
		switch e := err.(type) {
		case NameNotFoundError:
			log.Printf("HTTP %s - %d", e.Err, e.Code)
			http.Error(w, e.Err.Error(), e.Code)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func getName(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	name := vars["name"]
	if strings.EqualFold(name, "foo") {
		fmt.Fprintf(w, "Hello "+name)
		return nil
	} else {
		return NameNotFoundError{500, errors.New("name not found")}
	}
}

func main() {
	router := mux.NewRouter()
	router.Handle("/employee/get/{name}", WrapperHandler(getName)).Methods("GET")
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("Error starting http server: ", err)
		return
	}
}
