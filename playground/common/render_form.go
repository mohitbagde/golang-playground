package common

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type oauthForm struct {
	HashKey           string
	ExpectedSignature string
	ActualSignature   string
	BaseString        string
	SuccessMsg        string
}

// RenderOauth renders the Oauth form
func RenderOauth(w http.ResponseWriter, actualSignature string, oauth *OauthSignature, filepath string) {
	// Read the static error page, and render it
	t, err := template.ParseFiles(filepath)
	if err != nil {
		fmt.Println("Error in opening/parsing file ", err)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html")

	// Define the form to be displayed
	form := oauthForm{
		HashKey:           oauth.secret,
		ExpectedSignature: oauth.signature,
		ActualSignature:   actualSignature,
		BaseString:        oauth.baseString,
	}

	// Set the success message for the page
	form.SuccessMsg = "OAuth verification successful!"
	if !strings.EqualFold(actualSignature, oauth.signature) {
		form.SuccessMsg = "OAuth Failure! Signatures do not match!"
	}
	if err2 := t.Execute(w, form); err2 != nil {
		fmt.Println("Error in executing ", err2)
		return
	}
}
