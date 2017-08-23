package common

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"

	"math/rand"

	"strings"
	"time"
)

const nonceLength = 8
const oauthVersion = "1.0"
const oauthSignatureMethod = "HMAC-SHA1"

var alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// OauthSignature helpers
type OauthSignature struct {
	method      string
	scheme      string
	host        string
	path        string
	data        url.Values
	key         string
	secret      string
	tokenSecret string
	hashKey     string
	baseString  string
	signature   string
}

// NewOauthSignature Creates a new instance of OauthSignature
func NewOauthSignature(method, scheme, host, path string, data url.Values, key string, secret string) *OauthSignature {
	return &OauthSignature{method: method, scheme: scheme, host: host, path: path, data: data, key: key, secret: secret}
}

// CalcOAuthSignature Calculates an OAuth signature for the request
func (o *OauthSignature) CalcOAuthSignature(ctx context.Context) string {
	// Parse request of unnecessary OAuth parameters (not required to calculate signature)
	o.data.Del("oauth_signature")
	o.data.Del("oauth_token")
	// Generate the key from client/consumer secret & token secret in Router config
	//@TODO if a token secret is passed, use that as well for the signing key
	key := fmt.Sprintf("%s&", o.secret)

	// Prepare the base string from the remaining request parameters
	o.baseString = fmt.Sprintf(
		"%s&%s%s%s&%s",
		url.QueryEscape(o.method),
		url.QueryEscape(o.scheme),
		url.QueryEscape(o.host),
		url.QueryEscape(o.path),
		url.QueryEscape(strings.Replace(o.data.Encode(), "+", "%20", -1)))

	// Generate signature based on key and base string and return it
	signature, err := o.calcHash(key)
	if err != nil {
		fmt.Println("Error calculating hash")
	}
	o.signature = signature

	return o.signature
}

// Calculate the oauth_signature value based on the base string and the key
func (o *OauthSignature) calcHash(key string) (string, error) {
	mac := hmac.New(sha1.New, []byte(key))
	_, err := mac.Write([]byte(o.baseString))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), err
}

// Helper function to generate timestamp
func timestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

// Helper function to generate nonce
func nonce() string {
	nonce := make([]byte, nonceLength)
	for i := 0; i < nonceLength; i++ {
		nonce[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return string(nonce)
}

// Helper function to strip previous OAuth parameters
func (o *OauthSignature) stripOAuth() {
	for k, v := range o.data {
		for range v {
			if strings.Contains(strings.ToLower(k), "oauth") {
				o.data.Del(k)
			}
		}
	}
}

// PrepareOAuth Helper function to prepare New OAuth parameters
func (o *OauthSignature) PrepareOAuth(ctx context.Context) url.Values {
	o.stripOAuth()

	// Set new OAuth parameters
	o.data.Set("oauth_consumer_key", o.key)
	o.data.Set("oauth_signature_method", oauthSignatureMethod)
	o.data.Set("oauth_timestamp", timestamp())
	o.data.Set("oauth_nonce", nonce())
	o.data.Set("oauth_version", oauthVersion)
	o.data.Set("oauth_signature", o.CalcOAuthSignature(ctx))

	return o.data
}
