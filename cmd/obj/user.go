package obj

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"strings"
	"vobe-auth/cmd/db"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Username string
	Email    string
	Password string
	Scopes   []Scope
	Groups   []primitive.ObjectID

	tokens []Token
}

func (u User) Push() bool {
	_, ioErr := db.GetCollection("users").InsertOne(*db.CTX, u)
	if ioErr != nil {
		return false
	}
	return true
}

func (u User) GetTokens() []Token {

	if len(u.tokens) > 0 {
		return u.tokens
	}

	var _tokens []Token
	if findRes, findErr := db.GetCollection("tokens").Find(*db.CTX, bson.M{"user": u.ID}); findErr == nil {
		for findRes.Next(*db.CTX) {
			var t Token

			if decodeErr := findRes.Decode(&t); decodeErr != nil {
				fmt.Println("DECODE ERR\n", decodeErr)
				continue
			}

			_tokens = append(u.tokens, t)
		}

	}

	return _tokens
}
func FindUser(bson bson.M) (*User, error) {
	singleRes := db.GetCollection("users").FindOne(*db.CTX, bson)
	return WrapSingleRes(singleRes)
}

func WrapSingleRes(singleRes *mongo.SingleResult) (*User, error) {
	if singleRes == nil {
		return nil, errors.New("username not found")
	}
	user := &User{}
	if err := singleRes.Decode(user); err != nil {
		return nil, err
	}

	return user, nil
}

func hash(str string) string {
	input := strings.NewReader(str)

	hash := sha512.New()
	if _, err := io.Copy(hash, input); err != nil {
		log.Fatal(err)
	}

	return string(hash.Sum(nil))
}

func salt(str string) string {
	split := strings.Split(str, "")
	for i, val := range split {
		split[i] = string(([]byte(val)[0] << i * 2) % 0xff)
	}
	return strings.Join(split, "")
}

func Encrypt(pass string) string {
	return hash(salt(pass))
}
