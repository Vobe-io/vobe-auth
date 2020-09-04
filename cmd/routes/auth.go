package routes

import (
	. "vobe-auth/cmd/status"
)

func Auth(s string) *Status {
	return NewStatus(false, "Auth..")
}
