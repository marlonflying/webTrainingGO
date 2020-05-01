package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

//Route, in this way can be declared the routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes ... collection of routes
type Routes []Route

var routes = Routes{
	Route{
		"getEmployees",
		"GET",
		"/employees",
		getEmployees,
	},
	Route{
		"getEmployee",
		"GET",
		"/employee/{id}",
		getEmployee,
	}}

//Employee ... struct
type Employee struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//Employees ... collection of
type Employees []Employee

var employees []Employee

func init() {
	employees = Employees{
		Employee{Id: "1", FirstName: "Foo", LastName: "Bar"},
		Employee{Id: "2", FirstName: "Baz", LastName: "Qux"},
	}
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(employees)
}

func getEmployee(w http.ResponseWriter, r *http.Request) {
	vasrs := mux.Vars(r)
	id := vasrs["id"]
	for _, employee := range employees {
		if employee.Id == id {
			if err := json.NewEncoder(w).Encode(employee); err != nil {
				log.Print("error getting requested employee :: ", err)
			}
		}
	}
}

//AddRoutes ...
func AddRoutes(router *mux.Router) *mux.Router {
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

func main() {
	muxRouter := mux.NewRouter().StrictSlash(true)
	router := AddRoutes(muxRouter)
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}
}
