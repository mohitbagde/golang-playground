package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) init() {
	r := mux.NewRouter()

	// Static (UI) Handler
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir("./ui"))))

	sub := r.PathPrefix(fmt.Sprintf("/")).Subrouter()
	sub.Handle("/health", s.healthHandler(s.ctx)).Methods("GET")
	sub.Handle("/oauth", s.oauthHandler(s.ctx)).Methods("POST")

	s.ServeMux = http.NewServeMux()
	s.ServeMux.Handle("/", r)

}
