package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Server The web server instance
type Server struct {
}

// Create Create the web server
func Create() *Server {
	return &Server{}
}

// Listen Start listening
func (server *Server) Listen() error {
	var router = mux.NewRouter()
	router.HandleFunc("/media", media).Methods("GET")
	return http.ListenAndServe(":3000", router)
}

func media(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Still alive!")
}
