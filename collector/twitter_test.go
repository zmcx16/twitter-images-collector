package collector

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
		expectedToekn  string
		respStatusCode int
		respContent    string
		apiKey         string
		apiSecret      string
	}{
		{expectedToekn: "OOOXXX", respStatusCode: 200, respContent: "{\"token_type\":\"bearer\",\"access_token\":\"OOOXXX\"}", apiKey: "ABC", apiSecret: "DEF"},
		{expectedToekn: "", respStatusCode: 403, respContent: "", apiKey: "", apiSecret: ""},
	}

	for _, test := range tests {
		Convey(fmt.Sprintf("apiKey:\"%s\", apiSecret:\"%s\" => expectedToekn:\"%s\"", test.apiKey, test.apiSecret, test.expectedToekn), t, func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: test.respStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(test.respContent)),
					Header:     make(http.Header),
				}
			})

			twitterAPI := &TwitterAPI{Client: client}

			token := twitterAPI.GenBearerToken(test.apiKey, test.apiSecret)

			So(token, ShouldEqual, test.expectedToekn)
		})
	}
}

func TestGetTweets_GetTweetsData_ReturnTweets(t *testing.T) {

	tests := []struct {
		expected       bool
		respStatusCode int
		respContent    string
		token          string
		user           string
		start          string
		rts            bool
	}{
		{expected: true, respStatusCode: 200, respContent: `[{"entities":{"media":[]},"full_text":"OOOXXX"},{"entities":[{"media":[]}],"full_text":"XXXOOO"}]`, token: "XXX", user: "aaa", start: "0", rts: true},
		{expected: false, respStatusCode: 404, respContent: "[]", token: "XXX", user: "user not found", start: "0", rts: true},
	}

	for _, test := range tests {
		Convey(fmt.Sprintf("valid user: %t", test.expected), t, func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: test.respStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(test.respContent)),
					Header:     make(http.Header),
				}
			})

			twitterAPI := &TwitterAPI{Client: client}

			tweets := twitterAPI.GetTweets(test.token, test.user, test.start, test.rts)

			if test.expected {
				So(len(tweets), ShouldBeGreaterThanOrEqualTo, 0)
			} else {
				So(len(tweets), ShouldEqual, 0)
			}
		})
	}
}
