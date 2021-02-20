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

func genBearerToken(APIKey, APISecret string) string {

	credential := base64.StdEncoding.EncodeToString([]byte(APIKey + ":" + APISecret))
	data := url.Values{"grant_type": {"client_credentials"}}

	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Basic "+credential)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	// dump, _ := httputil.DumpRequestOut(req, true)
	// log.Println(string(dump))

	clt := http.Client{}
	resp, err := clt.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// dump, _ = httputil.DumpResponse(resp, true)
	// log.Println(string(dump))

	if resp.StatusCode == 200 {
		var jsonTokenResp BearerTokenResp
		json.Unmarshal(content, &jsonTokenResp)
		return jsonTokenResp.AccessToken
	}

	log.Fatal("Error! resp.StatusCode = " + strconv.Itoa(resp.StatusCode))
	return ""
}

func getTweets(token, user string, start string, rts bool) []map[string]interface{} {

	reqParam := "?screen_name=" + user + "&count=200&include_rts=" + strconv.FormatBool(rts) + "&tweet_mode=extended"
	if start != "0" {
		reqParam += "&max_id=" + start
	}

	req, err := http.NewRequest("GET", "https://api.twitter.com/1.1/statuses/user_timeline.json"+reqParam, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	clt := http.Client{}
	resp, err := clt.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(content))

	var jsonResp []map[string]interface{}
	if resp.StatusCode == 200 {
		json.Unmarshal(content, &jsonResp)
		return jsonResp
	}

	fmt.Println("Error! resp.StatusCode = " + strconv.Itoa(resp.StatusCode))
	return jsonResp
}
