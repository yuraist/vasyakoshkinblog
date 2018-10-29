package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/yuraist/vasyakoshkin/dblayer"
	"net/http"
)

type BlogHandler struct {
	dbHandler mongolayer.MongoLayer
}

func NewBlogHandler(dbHandler mongolayer.MongoLayer) *BlogHandler {
	return &BlogHandler{dbHandler:dbHandler,}
}

func HandleError(err error, whenCase string) {
	fmt.Printf("An error { %v } has been occured when { %v } is executing\n", err, whenCase)
}

func (handler *BlogHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	users, _ := handler.dbHandler.FindAll()

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err := json.NewEncoder(w).Encode(users)

	if err != nil {
		HandleError(err, "encoding user list")
	}
}

func (handler *BlogHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if username, ok := vars["username"]; ok {
		user, err := handler.dbHandler.FindUser(username)

		if err != nil {
			HandleError(err, "finding for a user")
		}

		w.Header().Set("Content-Type", "application/json;charset=utf8")
		json.NewEncoder(w).Encode(user)
	}
}

func (handler *BlogHandler) NewUser(w http.ResponseWriter, r *http.Request) {
	var user mongolayer.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("Error %v has been occured when json is decoding\n", err)
		return
	}

	err = handler.dbHandler.CreateUser(user)

	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("Error %v has been occured when user is adding to the database\n", err)
		return
	}

	success := mongolayer.SuccessResponse{Message: "user is added", Code: 201,}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(success)
}