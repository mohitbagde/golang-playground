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
	expectedSignature := "dMZy04/9T7Mku2zXs6sGDYmkctc="
	oauth := new(OauthSignature)
	// Define sample base string and key and MAC values
	baseString := "POST&https%3A%2F%2Flure-qastg.integration.nonprod.mheducation.com%2Fv1%2Fltiproxy&accept_media_types%3Dapplication%252Fvnd.ims.lti.v1.ltilink%26accept_multiple%3Dtrue%26accept_presentation_document_targets%3Diframe%26accept_unsigned%3Dfalse%26auto_create%3Dfalse%26can_confirm%3Dtrue%26content_item_return_url%3Dhttps%253A%252F%252Flogin-aws-qa.mhcampus.com%252FLTIHandlers%252FLTIContentItem.ashx%26context_id%3D561d7270e4b0279f288e8a70-57bc4628da25edde0464fc68d80d46282476919741bcdfb587fb7e35f0ffc47c%26context_label%3Dmcgmedium_91_qa10%26context_title%3DQA10%2520mcg%2520Medium%2520Class%2520000000091%26custom_correlation_id%3DReq05c03mk0i07ln3h7r8suz5%26custom_lms_type%3DBlackboard%26custom_mhe_org_id%3Dfe761c3e-6ddc-41e8-b063-a57c63d1fcff%26data%3D%257B%2522returnURL%2522%253A%2522https%253A%252F%252Fcsub2integrationqa10.blackboard.com%252Fwebapps%252Fpartner-cloud%252Fapp%252Fcontent%252Freturn%253Froot%253Dtrue%2526content_id%253D_39575_1%2526course_id%253D_7146_1%2526provider_id%253DMHE3%2526launch_time%253D1484941292623%2526launch_id%253D35205a2a-0b61-4775-b3a2-26854b7c325c%2526link_id%253Dca41ae808c4159d674462c260d2d6944%2522%252C%2522customerId%2522%253A%2522GBKB-YU5I-2RO4%2522%252C%2522lms%2522%253A%2522Blackboard%2522%252C%2522toolId%2522%253A%2522Connect2%2522%252C%2522originalRecordId%2522%253Anull%252C%2522contextId%2522%253A%2522561d7270e4b0279f288e8a70-57bc4628da25edde0464fc68d80d46282476919741bcdfb587fb7e35f0ffc47c%2522%257D%26launch_presentation_document_target%3Dframe%26launch_presentation_locale%3Den-US%26lis_person_contact_email_primary%3Dmcginstructor151_qa10%2540qatenemail.com%26lis_person_name_family%3DOneFiveOne%26lis_person_name_full%3Dqatenmcginstructor%2520OneFiveOne%26lis_person_name_given%3Dqatenmcginstructor%26lti_message_type%3DContentItemSelectionRequest%26lti_version%3DLTI-1p0%26oauth_consumer_key%3DMHCAMPUS%26oauth_nonce%3DNjM2MjA1MjAwOTQ2MjQ3NDY5%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1484941295%26oauth_version%3D1.0%26roles%3DInstructor%26tool_consumer_info_product_family_code%3DBlackboard%26tool_consumer_instance_description%3DBlackboard%2520cloud%2520services%2520for%2520partners%26tool_consumer_instance_guid%3DGBKB-YU5I-2RO4%26tool_consumer_instance_name%3DBb%2520Partner%2520Cloud%2520Service%26user_id%3Dmcginstructor151_qa10"
	key := "TW3L5RouyqS2zVz8&"

	// Assert equal
	hash, err := oauth.calcHash(baseString, key)
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
