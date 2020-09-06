package req

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"vobe-auth/cmd/db"
	"vobe-auth/cmd/obj"
)

type SignupReq struct {
	Username string
	Email    string
	Password string
}

type ISignupReq interface {
	exists(req SignupReq) bool
}

func (req SignupReq) Exists() bool {
	if count, _ := db.GetCollection("users").CountDocuments(*db.CTX, bson.M{"username": req.Username}); count > 0 {
		return true
	}
	return false
}

func (req SignupReq) ToUser() obj.User {
	return obj.User{
		ID:       primitive.NewObjectID(),
		Username: req.Username,
		Password: obj.Encrypt(req.Password),
		Email:    req.Email,
		Scopes:   []obj.Scope{},
		Groups:   []primitive.ObjectID{},
	}
}
