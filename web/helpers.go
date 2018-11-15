package web

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pauldotknopf/automounter/providers"
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

func convertMediaToMap(media []providers.Media) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	for _, media := range media {
		m := make(map[string]interface{})
		m["id"] = media.ID()
		m["displayName"] = media.DisplayName()
		m["provider"] = media.Provider()
		m["properties"] = media.Properties()
		result = append(result, m)
	}

	return result
}
