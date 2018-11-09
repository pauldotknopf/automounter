package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

	h := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: router}

	go func() {
		<-ctx.Done()
		fmt.Println("Shutting down the HTTP server...")
		h.Shutdown(ctx)
	}()

	err := h.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (server *Server) media(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	result := make([]map[string]interface{}, 0)

	for _, media := range server.mediaProvider.GetMedia() {
		m := make(map[string]interface{})
		m["id"] = media.ID()
		result = append(result, m)
	}

	j, _ := json.Marshal(result)
	io.WriteString(w, string(j))
}
