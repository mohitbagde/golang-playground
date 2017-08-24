package server

import (
	"context"
	"golang-playground/playground/common"
	"net/http"
	"path/filepath"
)

func (s *Server) healthHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(ny): CORS middleware with more restrictive settings.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		_, _ = w.Write([]byte("Health Check passed!"))
	})
}

func (s *Server) oauthHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte("Unable to parse Request!"))
			return
		}
		// Display the OAuth received from MHCampus
		actualMAC := r.PostForm.Get("oauth_signature")
		key := r.PostFormValue("oauth_consumer_key")
		oauth := common.NewOauthSignature(r.Method, scheme, r.Host, r.URL.Path, r.PostForm, key, "secret")

		// Verify that the OAuth signatures match
		oauth.CalcOAuthSignature(ctx)
		absPath, _ := filepath.Abs("./ui/index.html")
		common.RenderOauth(w, actualMAC, oauth, absPath)
	})
}
