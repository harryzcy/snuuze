package types

import "strconv"

type RequestFailedError struct {
	StatusCode int
}

func (e *RequestFailedError) Error() string {
	return "request failed with status code " + strconv.Itoa(e.StatusCode)
}
