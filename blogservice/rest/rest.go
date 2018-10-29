package rest

import (
	"github.com/gorilla/mux"
	"github.com/yuraist/vasyakoshkin/dblayer"
	"log"
	"net/http"
)

func ServeRestAPI(endpoint string, dbHandler mongolayer.MongoLayer) {
	handler := NewBlogHandler(dbHandler)
	router := mux.NewRouter()
	router.HandleFunc("/users", handler.ListAllUsers).Methods("GET")
	router.HandleFunc("/user/new", handler.NewUser).Methods("POST")
	router.HandleFunc("/user/{username}", handler.GetUser).Methods("GET")
	log.Fatal(http.ListenAndServe(endpoint, router))
}