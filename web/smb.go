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

type smbAddRequest struct {
	smbTestRequest
}

type smbAddResponse struct {
	genericResponse
	MediaID string `json:"mediaId"`
}

type smbRemoveRequest struct {
	MediaID string `json:"mediaId"`
}

type smbRemoveResponse struct {
	genericResponse
}

type smbDynamicLeaseRequest struct {
	smbTestRequest
}

type smbDynamicLeaseResponse struct {
	leaseCreateResponse
}

func (server *Server) smb(w http.ResponseWriter, r *http.Request) {
	var response smbResponse
	response.Success = true
	response.Entries = convertMediaArrayToJSON(server.smbProvider.GetMedia())
	sendResponse(w, http.StatusOK, response)
}

func (server *Server) smbTest(w http.ResponseWriter, r *http.Request) {

	var request smbTestRequest
	var response smbTestResponse

	err := getRequestBody(r, &request)
	if err != nil {
		sendError(w, err)
		return
	}

	// The request was a success (but maybe not the smb test)
	response.Success = true

	options, err := smb.CreateOptions(request.Server, request.Share, request.Folder, request.Security, request.Secure, request.Domain, request.Username, request.Password)
	if err != nil {
		sendError(w, err)
		return
	}

	err = server.smbProvider.TestConnection(options)
	if err != nil {
		response.Message = err.Error()
		response.IsValid = false
	} else {
		response.IsValid = true
	}

	sendResponse(w, http.StatusOK, response)
}

func (server *Server) smbAdd(w http.ResponseWriter, r *http.Request) {

	var request smbAddRequest
	var response smbAddResponse

	err := getRequestBody(r, &request)
	if err != nil {
		sendError(w, err)
		return
	}

	options, err := smb.CreateOptions(request.Server, request.Share, request.Folder, request.Security, request.Secure, request.Domain, request.Username, request.Password)
	if err != nil {
		response.Message = err.Error()
		response.Success = false
	} else {
		media, err := server.smbProvider.AddMedia(options)
		if err != nil {
			response.Message = err.Error()
			response.Success = false
		} else {
			response.MediaID = media.ID()
			response.Success = true
		}
	}

	sendResponse(w, http.StatusOK, response)
}

func (server *Server) smbRemove(w http.ResponseWriter, r *http.Request) {

	var request smbRemoveRequest
	var response smbRemoveResponse

	err := getRequestBody(r, &request)
	if err != nil {
		sendError(w, err)
		return
	}

	err = server.smbProvider.RemoveMedia(request.MediaID)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
	} else {
		response.Success = true
	}

	sendResponse(w, http.StatusOK, response)
}

func (server *Server) smbDynamicLease(w http.ResponseWriter, r *http.Request) {

	var request smbDynamicLeaseRequest
	var response smbDynamicLeaseResponse

	err := getRequestBody(r, &request)
	if err != nil {
		sendError(w, err)
		return
	}

	options, err := smb.CreateOptions(request.Server, request.Share, request.Folder, request.Security, request.Secure, request.Domain, request.Username, request.Password)
	if err != nil {
		response.Message = err.Error()
		response.Success = false
		sendResponse(w, http.StatusOK, response)
		return
	}

	// Build the media so that we can get the "id" to build the dynamic lease.
	lease, media, err := server.smbProvider.DynamicLease(options,
		server.leaser)

	if err != nil {
		response.Success = false
		response.Message = err.Error()
		sendResponse(w, http.StatusOK, response)
		return
	}

	response.Media = convertMediaToJSON(media)
	response.Success = true
	response.LeaseID = lease.ID()
	response.MountPath = lease.MountPath()

	sendResponse(w, http.StatusOK, response)
}
