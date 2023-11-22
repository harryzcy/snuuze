package types

import "strconv"

type RequestFailedError struct {
	For        string
	StatusCode int
	Body       string
}

func (e *RequestFailedError) Error() string {
	message := "request failed for " + e.For + " with status code " + strconv.Itoa(e.StatusCode)
	if e.Body != "" {
		message += ": " + e.Body
	}
	return message
}
