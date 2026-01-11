package scanner

import "fmt"

type opErr struct {
	pos     int
	msg     string
	content string
}

func (err opErr) Error() string {
	return fmt.Sprintf("%v at position %v; %v", err.msg, err.pos, err.content)
}
