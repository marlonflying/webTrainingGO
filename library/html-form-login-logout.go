package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

//User struct
type User struct {
	Username string
	Password string
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

func getUserName(request *http.Request) (userName string) {
	cookie, err := request.Cookie("session")
	if err == nil {
		cookieValue := make(map[string]string)
		err = cookieHandler.Decode("session", cookie.Value, &cookieValue)
		if err == nil {
			userName = cookieValue["username"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"username": userName,
	}
	encoded, err := cookieHandler.Encode("session", value)
	if err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func login(response http.ResponseWriter, request *http.Request) {

	user := readForm(request)
	target := "/"
	if user.Username != "" && user.Password != "" {
		setSession(user.Username, response)
		target = "/home"
	}
	http.Redirect(response, request, target, 302)
}

func logout(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("templates/login-form.html")
	parsedTemplate.Execute(w, nil)
}

func homePage(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		var data = map[string]interface{}{
			"userName": userName,
		}
		parsedTemplate, _ := template.ParseFiles("templates/home.html")
		parsedTemplate.Execute(response, data)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func main() {
	var router = mux.NewRouter()
	router.HandleFunc("/", loginPage)
	router.HandleFunc("/home", homePage)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout).Methods("POST")
	http.Handle("/", router)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("error starting http server: ", err)
		return
	}
}
