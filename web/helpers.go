package web

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func (server *Server) getRequestBody(r *http.Request, request interface{}) error {
	j, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, request)
}

func (server *Server) sendError(w http.ResponseWriter, err error) {
	var response genericResponse
	response.Success = false
	response.Message = err.Error()
	server.sendResponse(w, http.StatusBadRequest, response)
}

func (server *Server) sendResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(response)
	io.WriteString(w, string(j))
}
