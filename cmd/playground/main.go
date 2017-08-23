package main

import (
	"context"
	"golang-playground/playground/server"
)

func main() {
	ctx := context.Background()
	s := server.New(ctx)
	s.Start()
}
