package status

type Status struct {
	Success bool
	Code    int16
	Content string
}

func NewStatus(success bool, code int16, content string) Status {
	return Status{Success: success, Code: code, Content: content}
}
