package requestutil

import (
	"fmt"
	"io"
	"net/http"
)

func MustReadAll(resp *http.Response) []byte {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body", err)
	}
	return body
}
