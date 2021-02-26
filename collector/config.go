package collector

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// UserData struct
type UserData struct {
	UserID     string          `json:"user_id"`
	SaveDetail bool            `json:"save_detail"`
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

	twitterAPI ITwitterAPI
}

// LoadConfigFromPath (read config and generate twitter bearer token)
func (c *Config) LoadConfigFromPath(configPath string) bool {

	file, err := os.Open(configPath)
	if err != nil {
		log.Error(err)
	}
	defer file.Close()
	return c.LoadConfigFromReader(file)
}

// LoadConfigFromReader call by loadConfigFromPath for easy do unit test
func (c *Config) LoadConfigFromReader(r io.Reader) bool {

	byteConfig, err := ioutil.ReadAll(r)
	if err != nil {
		log.Error(err)
		return false
	}

	s := string(byteConfig)
	log.Println(s)

	json.Unmarshal(byteConfig, &c)

	if c.EnableLog {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	if c.twitterAPI == nil {
		c.twitterAPI = &TwitterAPI{Client: &http.Client{}}
	}

	c.BearerToken = c.twitterAPI.GenBearerToken(c.APIKey, c.APISecret)
	if c.BearerToken == "" {
		return false
	}

	return true
}
