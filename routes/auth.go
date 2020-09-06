package routes

import (
	"net/http"
	"vobe-auth/cmd/status"
)

func Auth(req *http.Request) status.Status {
	return status.NewStatus(false, 500, "Not implemented yet")
}
