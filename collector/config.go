package collector

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Config collector info and setting.
type Config struct {
	ThreadCnt    int      `json:"thread_cnt"`
	CollectUsers []string `json:"collect_users"`
	APIKey       string   `json:"api_key"`
	APISecret    string   `json:"api_secret"`
	BearerToken  string
}

// BearerTokenResp struct
type BearerTokenResp struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

// LoadConfig (read config and generate twitter bearer token)
func (c *Config) LoadConfig(configPath string) {

	byteConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(byteConfig, &c)
	fmt.Println(string(byteConfig))
	c.genToken()
	fmt.Println(c)
}

func (c *Config) genToken() {

	credential := base64.StdEncoding.EncodeToString([]byte(c.APIKey + ":" + c.APISecret))
	fmt.Println(credential)
	data := url.Values{"grant_type": {"client_credentials"}}

	req, err := http.NewRequest("POST", "https://api.twitter.com/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Basic "+credential)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	// dump, _ := httputil.DumpRequestOut(req, true)
	// fmt.Println(string(dump))

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
	// fmt.Println(string(dump))

	if resp.StatusCode == 200 {
		c.BearerToken = "ttt"

		var jsonTokenResp BearerTokenResp

		json.Unmarshal(content, &jsonTokenResp)
		c.BearerToken = jsonTokenResp.AccessToken
	}

}
