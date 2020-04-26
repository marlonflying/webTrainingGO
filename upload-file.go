package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println("Error getting the file from the form! : ", err)
		return
	}
	defer file.Close()
	out, pathError := os.Create("/tmp/uploadedFile")
	if pathError != nil {
		log.Println("Error creating the file: ", pathError)
		return
	}
	defer out.Close()
	_, copyFileError := io.Copy(out, file)
	if copyFileError != nil {
		log.Println("Error ocurred while file copy: ", copyFileError)
	}
	fmt.Fprintf(w, "File uploaded successfully: "+header.Filename)
}

func index(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("templates/upload-file.html")
	parsedTemplate.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", fileHandler)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("Error starting http server: ", err)
		return
	}
}
