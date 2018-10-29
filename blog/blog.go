package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

const (
	DB_NAME = "vasyakoshkinblog"
	DB_ADDRESS = "127.0.0.1"
	PORT = ":8000"
)

type User struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username"`
	FullName string `bson:"fullname" json:"fullname"`
	Age int `bson:"age" json:"age"`
}

func handleError(err error, whenCase string) {
	fmt.Printf("An error { %v } has been occured when { %v } is executing\n", err, whenCase)
}

func FindAll() ([]User, error) {
	session, err := mgo.Dial(DB_ADDRESS)
	var users = []User{}

	if err != nil {
		handleError(err, "connecting to the database")
		return nil, err
	}

	defer session.Close()

	err = session.DB(DB_NAME).C("users").Find(nil).All(&users)

	if err != nil {
		handleError(err, "users finding")
		return nil, err
	}

	return users, err
}

func FindUser(username string) (User, error) {
	session, err := mgo.Dial(DB_ADDRESS)
	user := User{}

	if err != nil {
		fmt.Printf("Error %v has been occured when connecting to the DB \n", err)
		return user, err
	}

	defer session.Close()

	collection := session.DB(DB_NAME).C("users")

	fmt.Printf("Searching for a %v user\n", username)
	err = collection.Find(bson.M{"username": username}).One(&user)

	if err != nil {
		fmt.Printf("Error %v has been occured when searching for a user\n", err)
		return user, err
	}

	return user, nil
}

func CreateUser(user User) (error) {
	session, err := mgo.Dial(DB_ADDRESS)

	if err != nil {
		return err
	}

	defer session.Close()

	collection := session.DB(DB_NAME).C("users")

	fmt.Printf("Creating ID for a %v user\n", user.Username)
	user.ID = bson.NewObjectId()

	fmt.Printf("Addind a %v user into the database\n", user.Username)
	err = collection.Insert(&user)

	if err != nil {
		return err
	}

	return nil
}

func ListAllUsers(w http.ResponseWriter, r *http.Request) {
	users, _ := FindAll()

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err := json.NewEncoder(w).Encode(users)

	if err != nil {
		handleError(err, "encoding user list")
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if username, ok := vars["username"]; ok {
		user, err := FindUser(username)

		if err != nil {
			handleError(err, "finding for a user")
		}

		w.Header().Set("Content-Type", "application/json;charset=utf8")
		json.NewEncoder(w).Encode(user)
	}
}

type SuccessResponse struct {
	message string
	code int
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("Error %v has been occured when json is decoding\n", err)
		return
	}

	err = CreateUser(user)

	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("Error %v has been occured when user is adding to the database\n", err)
		return
	}

	success := SuccessResponse{message: "user is added", code: 201,}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(success)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users", ListAllUsers).Methods("GET")
	router.HandleFunc("/user/new", NewUser).Methods("POST")
	router.HandleFunc("/user/{username}", GetUser).Methods("GET")
	log.Fatal(http.ListenAndServe(PORT, router))
}