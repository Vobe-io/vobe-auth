package status

type Status struct {
	Success bool
	Content string
}

func NewStatus(success bool, content string) *Status {
	return &Status{Success: success, Content: content}
}
