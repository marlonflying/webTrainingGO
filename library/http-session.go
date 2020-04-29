package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte("secret-key"))
}

func home(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	var autenticated interface{} = session.Values["authenticated"]
	if autenticated != nil {
		isAuthenticated := session.Values["authenticated"].(bool)
		if !isAuthenticated {
			http.Error(w, "You are not authorized to view the page!", http.StatusForbidden)
			return
		}
		fmt.Fprintln(w, "Home Page")
	} else {
		http.Error(w, "You are not authorized to view the page!", http.StatusForbidden)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Save(r, w)
	fmt.Fprintln(w, "You have successdully logged in!")
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	fmt.Fprintln(w, "You have logged out!")
}

func main() {
	http.HandleFunc("/home", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	err := http.ListenAndServe(connHost+":"+connPort, nil)
	if err != nil {
		log.Fatal("Error starting http server: ", err)
		return
	}
}
