package obj

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name string
	Scopes []Scope

}
