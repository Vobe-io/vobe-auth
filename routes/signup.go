package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"vobe-auth/cmd/db"
	"vobe-auth/cmd/obj/req"
	"vobe-auth/cmd/status"
)

func Signup(r *http.Request) status.Status {
	if body, err := ioutil.ReadAll(r.Body); err == nil {

		reqUser := req.SignupReq{}
		if err := json.Unmarshal(body, &reqUser); err == nil {

			if reqUser.Exists() {
				return status.NewStatus(false, 409, "User already exists")
			}

			user := reqUser.ToUser()
			if _, insertErr := db.GetCollection("users").InsertOne(*db.CTX, user); insertErr == nil {
				return status.NewStatus(true, 201, "Account created")
			} else {
				fmt.Println(insertErr)
				return status.NewStatus(false, 500, "Could not create account")
			}

		} else {
			return status.NewStatus(false, 500, "Can't parse body")
		}
	} else {
		return status.NewStatus(false, 500, "Can't read body")
	}
}
