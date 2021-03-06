package web

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pauldotknopf/automounter/leaser"
	"github.com/pauldotknopf/automounter/providers"
	"github.com/pauldotknopf/automounter/providers/smb"
)

// Server The web server instance
type Server struct {
	mediaProvider providers.MediaProvider
	leaser        leaser.Leaser
	smbProvider   smb.Provider
}

// Create Create the web server
func Create(leaser leaser.Leaser, smbProvider smb.Provider) *Server {
	return &Server{
		leaser.MediaProvider(),
		leaser,
		smbProvider,
	}
}

// Listen Start listening
func (server *Server) Listen(ctx context.Context, port int, started func()) error {
	var router = mux.NewRouter()
	router.HandleFunc("/media", server.media).Methods("GET")
	router.HandleFunc("/mount", server.mount).Methods("POST")
	router.HandleFunc("/unmount", server.unmount).Methods("POST")

	router.HandleFunc("/events", server.events)

	router.HandleFunc("/leases", server.leases).Methods("GET")
	router.HandleFunc("/leases/create", server.leaseCreate).Methods("POST")
	router.HandleFunc("/leases/release", server.leaseRelease).Methods("POST")

	if server.smbProvider != nil {
		router.HandleFunc("/smb", server.smb).Methods("GET")
		router.HandleFunc("/smb/test", server.smbTest).Methods("POST")
		router.HandleFunc("/smb/add", server.smbAdd).Methods("POST")
		router.HandleFunc("/smb/remove", server.smbRemove).Methods("POST")
		router.HandleFunc("/smb/dynamicLease", server.smbDynamicLease).Methods("POST")
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	h := &http.Server{Handler: router}

	go func() {
		<-ctx.Done()
		h.Shutdown(ctx)
	}()

	if started != nil {
		started()
	}

	err = h.Serve(l)
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (server *Server) media(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, convertMediaArrayToJSON(server.mediaProvider.GetMedia()))
}

func (server *Server) mount(w http.ResponseWriter, r *http.Request) {
	var request mountRequest
	getRequestBody(r, &request)

	if len(request.MediaID) == 0 {
		sendError(w, fmt.Errorf("no media id provided"))
		return
	}

	var response mountResponse

	session, err := server.mediaProvider.Mount(request.MediaID)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		sendResponse(w, http.StatusBadRequest, response)
		return
	}

	response.Success = true
	response.Location = session.Location()
	sendResponse(w, http.StatusOK, response)
}

func (server *Server) unmount(w http.ResponseWriter, r *http.Request) {
	var request unmountRequest
	getRequestBody(r, &request)

	if len(request.MediaID) == 0 {
		sendError(w, fmt.Errorf("no id provided"))
		return
	}

	var response unmountResponse

	err := server.mediaProvider.Unmount(request.MediaID)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		sendResponse(w, http.StatusBadRequest, response)
		return
	}

	response.Success = true
	sendResponse(w, http.StatusOK, response)
}
