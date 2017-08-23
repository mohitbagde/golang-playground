package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/facebookgo/httpdown"
)

// Server configuration
type Server struct {
	ServeMux *http.ServeMux
	ctx      context.Context
}

// New creates a server struct
func New(ctx context.Context) *Server {
	s := Server{
		ctx: ctx,
	}
	s.init()
	return &s
}

// Start the Server
func (s *Server) Start() {
	addr := fmt.Sprintf(":8003")

	server := &http.Server{
		Addr:    addr,
		Handler: s.ServeMux,
	}
	hd := &httpdown.HTTP{
		StopTimeout: 15,
		KillTimeout: 15,
	}

	fmt.Println("Listening on ", addr)
	if err := httpdown.ListenAndServe(server, hd); err != nil {
		panic(err)
	}
}
