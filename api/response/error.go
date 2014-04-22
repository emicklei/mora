package errhandler

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
	switch err.Error() {
	case "not found":
		code = http.StatusNotFound
	case "unauthorized":
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
