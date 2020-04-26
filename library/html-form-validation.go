package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
)

const (
	connHost          = "localhost"
	connPort          = "8080"
	userNameErrorMess = "PLease, enter a valid Username!"
	passErrorMess     = "Please, enter a valid Password!"
	genericErrorMess  = "Validation Error!"
)

//User struct
type User struct {
	Username string `valid:"alpha,required"`
	Password string `valid:"alpha,required"`
}

func readForm(r *http.Request) *User {
	r.ParseForm()
	user := new(User)
	decoder := schema.NewDecoder()
	decodeErr := decoder.Decode(user, r.PostForm)
	if decodeErr != nil {
		log.Println("error mapping parsed form data to struct : ", decodeErr)
	}
	return user
}

func validateUser(w http.ResponseWriter, r *http.Request, user *User) (bool, string) {
	valid, validationError := govalidator.ValidateStruct(user)
	if !valid {
		userNameError := govalidator.ErrorByField(validationError, "Username")
		passwordError := govalidator.ErrorByField(validationError, "Password")
		if userNameError != "" {
			log.Printf("username validation error: ", userNameError)
			return valid, userNameErrorMess
		}
		if passwordError != "" {
			log.Printf("Password validation error:", passwordError)
		}
	}
	return valid, genericErrorMess
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		parsedTemplate, _ := template.ParseFiles("templates/login-form.html")
		parsedTemplate.Execute(w, nil)
	} else {
		user := readForm(r)
		valid, validationErrorMessage := validateUser(w, r, user)
		if !valid {
			fmt.Fprintf(w, validationErrorMessage)
			return
		}
		fmt.Fprintf(w, "Validation Succesful, hello: "+user.Username+"!!!")
	}
}

func main() {
	http.HandleFunc("/", login)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}
