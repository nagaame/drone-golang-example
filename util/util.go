package util

import (
	"fmt"
	"net/http"
	"strings"
)

func SayHello(response http.ResponseWriter, r *http.Request) {
	message := r.URL.Path

	message = strings.TrimPrefix(message, "/")

	message = "Hello " + message + " : " + fmt.Sprintf("%d", GetNameLen(message))

	response.Write([]byte(message))
}

func GetNameLen(s string) int {
	return len(s)
}