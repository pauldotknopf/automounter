package web

import (
	"fmt"
	"net/http"
)

type leasesResponse struct {
	genericResponse
	Leases []map[string]interface{} `json:"leases"`
}

type leaseCreateRequest struct {
	MediaID string `json:"mediaId"`
}

type leaseCreateResponse struct {
	genericResponse
	MountPath string `json:"mountPath"`
	LeaseID   string `json:"leaseId"`
}

type leaseReleaseRequest struct {
	LeaseID string `json:"leaseId"`
}

type leaseReleaseResponse struct {
	genericResponse
}

func (server *Server) leases(w http.ResponseWriter, r *http.Request) {
	var response leasesResponse

	leases := server.leaser.Leases()
	response.Leases = make([]map[string]interface{}, 0)
	for _, lease := range leases {
		l := make(map[string]interface{})
		l["leaseId"] = lease.ID()
		l["mediaId"] = lease.MediaID()
		l["mountPath"] = lease.MountPath()
		l["isValid"] = lease.IsValid()
		response.Leases = append(response.Leases, l)
	}

	response.Success = true
	sendResponse(w, http.StatusOK, response)
}

func (server *Server) leaseCreate(w http.ResponseWriter, r *http.Request) {
	var request leaseCreateRequest
	getRequestBody(r, &request)

	if len(request.MediaID) == 0 {
		sendError(w, fmt.Errorf("no media id provided"))
		return
	}

	var response leaseCreateResponse

	lease, err := server.leaser.Lease(request.MediaID)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		sendResponse(w, http.StatusBadRequest, response)
		return
	}

	response.Success = true
	response.LeaseID = lease.ID()
	response.MountPath = lease.MountPath()
	sendResponse(w, http.StatusOK, response)
}

func (server *Server) leaseRelease(w http.ResponseWriter, r *http.Request) {
	var request leaseReleaseRequest
	getRequestBody(r, &request)

	if len(request.LeaseID) == 0 {
		sendError(w, fmt.Errorf("no lease id provided"))
		return
	}

	var response leaseReleaseResponse

	err := server.leaser.Release(request.LeaseID)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		sendResponse(w, http.StatusBadRequest, response)
		return
	}

	response.Success = true
	sendResponse(w, http.StatusOK, response)
}
