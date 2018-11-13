package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pauldotknopf/automounter/providers"
)

// Server The web server instance
type Server struct {
	mediaProvider providers.MediaProvider
}

// Create Create the web server
func Create(mediaProvider providers.MediaProvider) *Server {
	return &Server{
		mediaProvider,
	}
}

// Listen Start listening
func (server *Server) Listen(ctx context.Context, port int) error {
	var router = mux.NewRouter()
	router.HandleFunc("/media", server.media).Methods("GET")
	router.HandleFunc("/mount", server.mount).Methods("POST")
	router.HandleFunc("/unmount", server.unmount).Methods("POST")

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
	result := make([]map[string]interface{}, 0)

	for _, media := range server.mediaProvider.GetMedia() {
		m := make(map[string]interface{})
		m["id"] = media.ID()
		m["displayName"] = media.DisplayName()
		result = append(result, m)
	}

	server.sendResponse(w, http.StatusBadRequest, result)
}

func (server *Server) mount(w http.ResponseWriter, r *http.Request) {
	var request mountRequest
	server.getRequestBody(r, &request)

	if len(request.MediaID) == 0 {
		server.sendError(w, fmt.Errorf("no id provided"))
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
