package web

import (
	"net/http"

	"github.com/pauldotknopf/automounter/providers/smb"
)

type smbResponse struct {
	genericResponse
	Entries []map[string]interface{} `json:"entries"`
}

type smbTestRequest struct {
	Server   string `json:"server"`
	Share    string `json:"share"`
	Folder   string `json:"folder"`
	Security string `json:"security"`
	Secure   bool   `json:"secure"`
	Domain   string `json:"domain"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type smbTestResponse struct {
	genericResponse
	IsValid bool `json:"isValid"`
}

func (server *Server) smb(w http.ResponseWriter, r *http.Request) {
	var response smbResponse
	response.Success = true
	response.Entries = convertMediaToMap(server.smbProvider.GetMedia())
	server.sendResponse(w, http.StatusOK, response)
}

func (server *Server) smbTest(w http.ResponseWriter, r *http.Request) {

	var request smbTestRequest
	var response smbTestResponse

	err := server.getRequestBody(r, &request)
	if err != nil {
		server.sendError(w, err)
		return
	}

	// The request was a success (but maybe not the smb test)
	response.Success = true

	options, err := smb.CreateOptions(request.Server, request.Share, request.Folder, request.Security, request.Secure, request.Domain, request.Username, request.Password)
	if err != nil {
		server.sendError(w, err)
		return
	}

	err = server.smbProvider.TestConnection(options)
	if err != nil {
		response.Message = err.Error()
		response.IsValid = false
	} else {
		response.IsValid = true
	}

	server.sendResponse(w, http.StatusOK, response)
}
