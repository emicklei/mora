package response

import (
	"github.com/emicklei/go-restful"
	"log"
	"net/http"
)

type Error struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

func WriteError(err error, resp *restful.Response) {
	log.Printf("[mora][error] %v", err)

	// Set response status code
	code := http.StatusInternalServerError

	// String error
	error := err.Error()

	if error == "not found" || len(error) > 7 && error[:7] == "Unknown" {
		code = http.StatusNotFound
	} else if error == "unauthorized" || len(error) > 14 && error[:14] == "not authorized" {
		code = http.StatusUnauthorized
	}

	// Write error response
	WriteStatusError(code, err, resp)
}

func WriteStatusError(status int, err error, resp *restful.Response) {
	success := NewResponse(false)
	success.SetError(err)
	success.WriteStatus(status, resp)
}
