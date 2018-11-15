package web

import (
	"context"
	"fmt"
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
func (server *Server) Listen(ctx context.Context, port int) error {
	var router = mux.NewRouter()
	router.HandleFunc("/media", server.media).Methods("GET")
	router.HandleFunc("/mount", server.mount).Methods("POST")
	router.HandleFunc("/unmount", server.unmount).Methods("POST")

	router.HandleFunc("/leases", server.leases).Methods("GET")
	router.HandleFunc("/leases/create", server.leaseCreate).Methods("POST")
	router.HandleFunc("/leases/release", server.leaseRelease).Methods("POST")

	if server.smbProvider != nil {
		router.HandleFunc("/smb", server.smb).Methods("GET")
		router.HandleFunc("/smb/test", server.smbTest).Methods("POST")
		router.HandleFunc("/smb/add", server.smbAdd).Methods("POST")
	}

	h := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: router}

	go func() {
		<-ctx.Done()
		h.Shutdown(ctx)
	}()

	err := h.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (server *Server) media(w http.ResponseWriter, r *http.Request) {
	server.sendResponse(w, http.StatusOK, convertMediaToMap(server.mediaProvider.GetMedia()))
}

func (server *Server) mount(w http.ResponseWriter, r *http.Request) {
	var request mountRequest
	server.getRequestBody(r, &request)

	if len(request.MediaID) == 0 {
		server.sendError(w, fmt.Errorf("no media id provided"))
	}

	var response mountResponse

	session, err := server.mediaProvider.Mount(request.MediaID)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		server.sendResponse(w, http.StatusBadRequest, response)
		return
	}

	response.Success = true
	response.Location = session.Location()
	server.sendResponse(w, http.StatusOK, response)
}

func (server *Server) unmount(w http.ResponseWriter, r *http.Request) {
	var request unmountRequest
	server.getRequestBody(r, &request)

	if len(request.MediaID) == 0 {
		server.sendError(w, fmt.Errorf("no id provided"))
	}

	var response unmountResponse

	err := server.mediaProvider.Unmount(request.MediaID)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		server.sendResponse(w, http.StatusBadRequest, response)
		return
	}

	response.Success = true
	server.sendResponse(w, http.StatusOK, response)
}
