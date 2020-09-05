package catchers

import "vobe-auth/cmd/status"

func NotFound(msg string) status.Status {
	return status.NewStatus(false, 404, msg)
}
