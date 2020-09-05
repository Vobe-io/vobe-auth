package catchers

import "vobe-auth/cmd/status"

func InternalServerError(msg string) status.Status {
	return status.NewStatus(false, 500, msg)
}
