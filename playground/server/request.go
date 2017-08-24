package server

// ReceiveLTI parses lti request
import (
	"context"
	"errors"
	"fmt"
	"golang-playground/playground/common"
	"net/http"
	"strings"
)

const scheme = "http://"

// VerifyOAuth is a method that receives a request and validates oauth
func VerifyOAuth(ctx context.Context, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	// Display the OAuth received from MHCampus
	actualMAC := r.PostForm.Get("oauth_signature")
	key := r.PostFormValue("oauth_consumer_key")
	sig := common.NewOauthSignature(r.Method, scheme, r.Host, r.URL.Path, r.PostForm, key, "secret")

	// Verify that the OAuth signatures match
	if !strings.EqualFold(actualMAC, sig.CalcOAuthSignature(ctx)) {
		errOauth := fmt.Sprintf("OAuth Verification Failure")
		return errors.New(errOauth)
	}
	fmt.Println("OAuth Verification Successful")
	return nil
}
