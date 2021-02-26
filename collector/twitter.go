package collector

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// BearerTokenResp struct
type BearerTokenResp struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

// TwitterAPI struct
type TwitterAPI struct {
	Client *http.Client
}

// ITwitterAPI interface
type ITwitterAPI interface {
  GenBearerToken(APIKey, APISecret string) string
  GetTweets(token, user, start string, rts bool) []map[string]interface{}
}

// GenBearerToken generate twitter bearer token
func (tapi *TwitterAPI) GenBearerToken(APIKey, APISecret string) string {

	credential := base64.StdEncoding.EncodeToString([]byte(APIKey + ":" + APISecret))
	data := url.Values{"grant_type": {"client_credentials"}}

	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		log.Error(err)
	}
	req.Header.Add("Authorization", "Basic "+credential)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	// dump, _ := httputil.DumpRequestOut(req, true)
	// log.Println(string(dump))

	resp, err := tapi.Client.Do(req)
	if err != nil {
		log.Error(err)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	// dump, _ = httputil.DumpResponse(resp, true)
	// log.Println(string(dump))

	if resp.StatusCode == 200 {
		var jsonTokenResp BearerTokenResp
		json.Unmarshal(content, &jsonTokenResp)
		return jsonTokenResp.AccessToken
	}

	errorMsg := "Error! resp.StatusCode = " + strconv.Itoa(resp.StatusCode)
	log.Error(errorMsg)
	fmt.Println(errorMsg)
	return ""
}

// GetTweets get user timeline twitter
func (tapi *TwitterAPI) GetTweets(token, user, start string, rts bool) []map[string]interface{} {

	reqParam := "?screen_name=" + user + "&count=200&include_rts=" + strconv.FormatBool(rts) + "&tweet_mode=extended"
	if start != "0" {
		reqParam += "&max_id=" + start
	}

	req, err := http.NewRequest("GET", "https://api.twitter.com/1.1/statuses/user_timeline.json"+reqParam, nil)
	if err != nil {
		log.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := tapi.Client.Do(req)
	if err != nil {
		log.Error(err)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	// fmt.Println(string(content))

	var jsonResp []map[string]interface{}
	if resp.StatusCode == 200 {
		json.Unmarshal(content, &jsonResp)
		return jsonResp
	}

	errorMsg := "Error! resp.StatusCode = " + strconv.Itoa(resp.StatusCode)
	log.Error(errorMsg)
	fmt.Println(errorMsg)
	return jsonResp
}
