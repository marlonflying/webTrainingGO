package main

import (
	"fmt"
	"log"
	"net/http"

	redisStore "gopkg.in/boj/redistore.v1"
)

const (
	connHost = "localhost"
	connPort = "8080"
)

var store *redisStore.RediStore
var err error

func init() {
	store, err = redisStore.NewRediStore(10, "tcp", ":6379", "", []byte("secret-key"))
	if err != nil {
		log.Fatal("Error getting Redis Store: ", err)
	}
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
	if err = session.Save(r, w); err != nil {
		log.Fatalf("Error saving session: %v", err)
	}
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
	defer store.Close()
	if err != nil {
		log.Fatal("Error starting http server: ", err)
		return
	}
}
