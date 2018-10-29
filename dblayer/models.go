package mongolayer

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username"`
	FullName string `bson:"fullname" json:"fullname"`
	Age int `bson:"age" json:"age"`
}

type SuccessResponse struct {
	Message string
	Code int
}
