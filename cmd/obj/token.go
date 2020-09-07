package obj

import (
	"errors"
	"fmt"
	"time"
	"vobe-auth/cmd/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Token struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	User       primitive.ObjectID `bson:"user" json:"user"`
	Expiration int64              `bson:"expiration" json:"expiration"`
}

func (t Token) IsValid(u User) bool {
	return u.ID == t.ID && !t.IsExpired()
}

func (t Token) IsExpired() bool {
	return time.Unix(t.Expiration, 0).After(time.Now())
}

func (t Token) AsString() string {
	return t.ID.Hex()
}

func (t Token) Push() bool {
	_, ioErr := db.GetCollection("tokens").InsertOne(*db.CTX, t)
	if ioErr != nil {
		fmt.Println(ioErr)
		return false
	}
	return true
}

func NewDefaultToken(user User) Token {
	return NewToken(user, time.Hour*24*30)
}

func NewToken(user User, duration time.Duration) Token {

	return Token{
		ID:         primitive.NewObjectID(),
		Expiration: time.Now().Add(duration).Unix(),
		User:       user.ID,
	}
}

func GetFromString(token string) (*Token, error) {
	if token, err := WrapSingleTokenRes(db.GetCollection("tokens").FindOne(*db.CTX, bson.M{"_id": token})); err != nil {
		return nil, errors.New("token not found")
	} else {
		return token, nil
	}
}

func WrapSingleTokenRes(singleRes *mongo.SingleResult) (*Token, error) {
	if singleRes == nil {
		return nil, errors.New("token not found")
	}
	token := &Token{}
	if err := singleRes.Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}
