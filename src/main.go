package main

import (
	"context"
	"net/http"

	"github.com/facebookgo/httpdown"
)

// Server configuration
type Server struct {
	ServeMux *http.ServeMux
	ctx      context.Context
}

func main() {

	ctx := context.Background()
	s := Server{
		ctx: ctx,
	}
	s.ServeMux = http.NewServeMux()

	server := &http.Server{
		Addr:    "8001",
		Handler: s.ServeMux,
	}
	hd := &httpdown.HTTP{
		StopTimeout: 15,
		KillTimeout: 15,
	}

	if err := httpdown.ListenAndServe(server, hd); err != nil {
		panic(err)
	}
}
