package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

//Person type
type Person struct {
	Id   string
	Name string
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	person := Person{Id: "44", Name: "Careloco"}
	parsedTemplate, _ := template.ParseFiles("templates/first-template.html")
	err := parsedTemplate.Execute(w, person)
	if err != nil {
		log.Println("Error occurred while rendering the template or writing it's output : ", err)
		return
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", renderTemplate).Methods("GET")
	router.PathPrefix("/").Handler(http.StripPrefix("/static",
		http.FileServer(http.Dir("static/"))))

	http.FileServer(http.Dir("static/"))
	err := http.ListenAndServe(connHost+":"+connPort, router)
	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}
