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
