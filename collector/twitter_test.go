package collector

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestGenBearerToken_GenerateToken_ReturnToken(t *testing.T) {

	tests := []struct {
		respStatusCode int
		respContent    string
		apiKey         string
		apiSecret      string
	}{
		{respStatusCode: 200, respContent: "{\"token_type\":\"bearer\",\"access_token\":\"OOOXXX\"}", apiKey: "ABC", apiSecret: "DEF"},
		{respStatusCode: 403, respContent: "", apiKey: "", apiSecret: ""},
	}

	for _, test := range tests {
		client := NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: test.respStatusCode,
				Body:       ioutil.NopCloser(bytes.NewBufferString(test.respContent)),
				Header:     make(http.Header),
			}
		})

		twitterAPI := &TwitterAPI{Client: client}

		token := twitterAPI.GenBearerToken(test.apiKey, test.apiSecret)
		if test.apiKey != "" && test.apiSecret != "" {
			assert.NotEmpty(t, token, test)
		} else {
			assert.Equal(t, "", token, test)
		}
	}
}

func TestGetTweets_GetTweetsData_ReturnTweets(t *testing.T) {

	tests := []struct {
		respStatusCode 	int
		respContent    	string
		token         	string
		user      			string
		start						string
		rts							bool
	}{
		{respStatusCode: 200, respContent: "[{\"entities\":[{\"media\":[]}],\"full_text\":\"OOOXXX\"},{\"entities\":[{\"media\":[]}],\"full_text\":\"XXXOOO\"}]", token: "XXX", user: "aaa", start: "0", rts: true},
		{respStatusCode: 404, respContent: "[]", token: "XXX", user: "user not found", start: "0", rts: true},
	}

	for _, test := range tests {
		client := NewTestClient(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: test.respStatusCode,
				Body:       ioutil.NopCloser(bytes.NewBufferString(test.respContent)),
				Header:     make(http.Header),
			}
		})

		twitterAPI := &TwitterAPI{Client: client}

		tweets := twitterAPI.GetTweets(test.token, test.user, test.start, test.rts)
		if test.user != "user not found" {
			assert.True(t, len(tweets)>0, test)
		} else {
			assert.True(t, len(tweets)==0, test)
		}
	}
}

