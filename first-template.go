package main

import (
	"html/template"
	"log"
	"net/http"
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
	person := Person{Id: "23", Name: "Careloco"}
	parsedTemplate, _ := template.ParseFiles("templates/first-template.html")
	err := parsedTemplate.Execute(w, person)
	if err != nil {
		log.Println("Error occurred while rendering the template or writing it's output : ", err)
		return
	}
}

func main() {
	http.HandleFunc("/", renderTemplate)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}
