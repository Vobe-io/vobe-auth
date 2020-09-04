package catchers

import (
	. "vobe-auth/cmd/status"
)

func NotFound(s string) *Status {
	return NewStatus(false, "Not found")
}
