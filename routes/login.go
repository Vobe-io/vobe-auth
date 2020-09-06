package routes

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
	"vobe-auth/cmd/obj"
	"vobe-auth/cmd/obj/req"
	"vobe-auth/cmd/status"
)

func Login(r *http.Request) status.Status {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return status.NewStatus(false, 500, "Can't read body")
	}

	userData := req.LoginReq{}
	err = json.Unmarshal(body, &userData)

	if err != nil {
		return status.NewStatus(false, 500, "Can't parse body to json")
	}

	if user, err := obj.FindUser(bson.M{"$or": []bson.M{{"username": userData.Username}, {"email": userData.Username}}}); err != nil {
		return status.NewStatus(false, 404, "User not found")
	} else {
		if user.Password == obj.Encrypt(userData.Password) {
			token := obj.NewDefaultToken(*user)
			if token.Push() {
				return status.NewStatus(true, 200, token.AsString())
			}

		}
		return status.NewStatus(false, 500, "idk man")
	}
}
