package web

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pauldotknopf/automounter/providers"
)

func getRequestBody(r *http.Request, request interface{}) error {
	j, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, request)
}

func sendError(w http.ResponseWriter, err error) {
	var response genericResponse
	response.Success = false
	response.Message = err.Error()
	sendResponse(w, http.StatusBadRequest, response)
}

func sendResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	j, _ := json.Marshal(response)
	io.WriteString(w, string(j))
}

func convertMediaToJSON(media providers.Media) map[string]interface{} {
	m := make(map[string]interface{})
	m["id"] = media.ID()
	m["displayName"] = media.DisplayName()
	m["provider"] = media.Provider()
	m["properties"] = media.Properties()
	return m
}

func convertMediaArrayToJSON(media []providers.Media) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, media := range media {
		result = append(result, convertMediaToJSON(media))
	}
	return result
}
