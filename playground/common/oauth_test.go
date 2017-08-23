package common

import (
	"net/http"
	"strconv"
	"time"

	"github.mheducation.com/MHEducation/dle-orchestration/orchestration/test"

	"github.com/stretchr/testify/assert"

	"net/url"
	"strings"
	"testing"
)

func TestCalcHash(t *testing.T) {
	assert := assert.New(t)
	expectedSignature := "zoARStOzWrXTl2liiTcpXqj57co="
	oauth := new(OauthSignature)
	// Define sample base string and key and MAC values
	oauth.baseString =
		`POST&http%3A%2F%2Flocalhost%3A8003%2F
	oauth&hello%3Dworld%26
	oauth_consumer_key%3Dkey%26
	oauth_nonce%3DHTJTlH%26
	oauth_signature_method%3DHMAC-SHA1%26
	oauth_timestamp%3D1503519169%26
	oauth_version%3D1.0`
	key := "TW3L5RouyqS2zVz8&"

	// Assert equal
	hash, err := oauth.calcHash(key)
	assert.NoError(err)
	assert.Equal(hash, expectedSignature)
}

func TestStripOAuth(t *testing.T) {
	assert := assert.New(t)
	method := "POST"
	scheme := "https://"
	host := "localhost:8000"
	path := "/ltiproxy"
	data := url.Values{
		"oauth_version":      {"1.0"},
		"oauth_consumer_key": {"test-key"},
		"lti_message_type":   {"basic-lti-launch-request"},
	}
	uri := scheme + host + path
	oauth := NewOauthSignature(method, scheme, host, path, data, "test-key", "test-secret")
	req, err := http.NewRequest(method, uri, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	assert.Nil(err)

	// Parse the request and assert data exists before stripping
	parseErr := req.ParseForm()
	assert.NoError(parseErr)
	assert.Equal("test-key", req.PostFormValue("oauth_consumer_key"))

	// Assert data has been stripped
	oauth.stripOAuth()
	req.PostForm = oauth.data
	assert.Equal("", req.PostFormValue("oauth_consumer_key"))
}

func TestPrepareOAuth(t *testing.T) {
	assert := assert.New(t)
	method := "POST"
	scheme := "https://"
	host := "localhost:8000"
	path := "/ltiproxy"
	data := url.Values{
		"lti_message_type": {"basic-lti-launch-request"},
	}
	uri := scheme + host + path
	ctx := test.NewTestContext()
	oauth := NewOauthSignature(method, scheme, host, path, data, "test-key", "test-secret")

	req, err := http.NewRequest(method, uri, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	assert.Nil(err)
	if err = req.ParseForm(); err != nil {
		assert.NotNil(err)
	}

	oauth.PrepareOAuth(ctx)
	req.PostForm = oauth.data

	assert.Equal("test-key", req.PostFormValue("oauth_consumer_key"))
	assert.Equal(oauthSignatureMethod, req.PostFormValue("oauth_signature_method"))
	assert.Equal(oauthVersion, req.PostFormValue("oauth_version"))
	oauthTimestamp, err := strconv.Atoi(req.PostFormValue("oauth_timestamp"))
	assert.NoError(err)
	assert.InDelta(time.Now().Unix(), oauthTimestamp, 0.05) // ensure oauth timestamp is within %5 of now
	assert.Nil(err)
}
