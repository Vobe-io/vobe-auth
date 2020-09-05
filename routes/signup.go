package routes

import (
	"net/http"
	"vobe-auth/cmd/status"
)

func Signup(r *http.Request) status.Status {

	return status.NewStatus(true, 200, "Signup complete")
}
