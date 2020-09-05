package routes

import (
	"net/http"
	"vobe-auth/cmd/status"
)

func Login(r *http.Request) status.Status {

	return status.NewStatus(true, 200, "Login complete")
}
