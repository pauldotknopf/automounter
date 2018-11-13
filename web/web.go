package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
