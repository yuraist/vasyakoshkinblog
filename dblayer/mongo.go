package mongolayer

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


const (
	DB         = "vasyakoshkinblog"
	CONNECTION = "127.0.0.1"
	USERS      = "users"
)

type MongoLayer struct {
	session *mgo.Session
}

func NewMongoLayer() (*MongoLayer, error) {
	s, err := mgo.Dial(CONNECTION)

	return &MongoLayer{s, }, err
}

func (db *MongoLayer) getFreshSession() *mgo.Session {
	return db.session.Copy()
}

func handleError(err error, whenCase string) {
	fmt.Printf("An { %v } error occured when { %v } is executing\n", err, whenCase)
}

func (db *MongoLayer) FindAll() ([]User, error) {
	session := db.getFreshSession()
	var users = []User{}

	defer session.Close()

	err := session.DB(DB).C(USERS).Find(nil).All(&users)

	if err != nil {
		handleError(err, "FINDING FOR ALL USERS")
		return nil, err
	}

	return users, err
}

func (db *MongoLayer) FindUser(username string) (User, error) {
	session := db.getFreshSession()
	user := User{}

	defer session.Close()

	collection := session.DB(DB).C(USERS)

	fmt.Printf("Searching for a %v user\n", username)
	err := collection.Find(bson.M{"username": username}).One(&user)

	if err != nil {
		handleError(err, "FINDING FOR A USER IN A DATABASE")
		return user, err
	}

	return user, nil
}

func (db *MongoLayer) CreateUser(user User) (error) {
	session := db.getFreshSession()

	defer session.Close()

	collection := session.DB(DB).C(USERS)

	fmt.Printf("Creating ID for a %v user\n", user.Username)
	user.ID = bson.NewObjectId()

	fmt.Printf("Addind a %v user into the database\n", user.Username)
	err := collection.Insert(&user)

	if err != nil {
		handleError(err, "ADDING A USER INTO A DATABASE")
		return err
	}

	return nil
}