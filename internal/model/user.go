package model

type User struct {
	Id       string `bson:"_id,omitempty" json:"id"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password,omitempty" json:"-"`
}
