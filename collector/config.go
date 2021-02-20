package collector

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// UserData struct
type UserData struct {
	UserID     string `json:"user_id"`
	FolderName string `json:"folder_name"`
	DestPath   string `json:"dest_path"`
}

// Config collector info and setting.
type Config struct {
	ThreadCnt       int        `json:"thread_cnt"`
	CollectUsers    []UserData `json:"collect_users"`
	ImageSize       string     `json:"image_size"`
	EnableLog       bool       `json:"enable_log"`
	SyncLastNDays   int        `json:"sync_last_n_days"`
	IncludeRetweets bool       `json:"include_retweets"`
	APIKey          string     `json:"api_key"`
	APISecret       string     `json:"api_secret"`
	BearerToken     string
}

// LoadConfig (read config and generate twitter bearer token)
func (c *Config) LoadConfig(configPath string) {

	byteConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(byteConfig, &c)

	if c.EnableLog {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	c.BearerToken = genBearerToken(c.APIKey, c.APISecret)
}

func (c *Config) genToken() {

	credential := base64.StdEncoding.EncodeToString([]byte(c.APIKey + ":" + c.APISecret))
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
		c.BearerToken = jsonTokenResp.AccessToken
	} else {
		log.Fatal("Error! resp.StatusCode = " + strconv.Itoa(resp.StatusCode))
	}

}
