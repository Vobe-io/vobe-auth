package obj

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	_id primitive.ObjectID
	username string
	password string
	scopes []Scope
	groups []primitive.ObjectID
}