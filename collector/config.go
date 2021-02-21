package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// UserData struct
type UserData struct {
	UserID     string          `json:"user_id"`
	FolderName json.RawMessage `json:"folder_name"`
	DestPath   json.RawMessage `json:"dest_path"`
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
func (c *Config) LoadConfig(configPath string) bool {

	byteConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Error(err)
		return false
	}

	json.Unmarshal(byteConfig, &c)

	if c.EnableLog {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	twitterAPI := &TwitterAPI{Client: &http.Client{}}
	c.BearerToken = twitterAPI.GenBearerToken(c.APIKey, c.APISecret)
	if c.BearerToken == "" {
		return false
	}

	return true
}
