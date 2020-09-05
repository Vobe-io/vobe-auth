package obj

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	_id primitive.ObjectID
	name string
	scopes []Scope
}
