package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// Collector struct
type Collector struct {
	config Config
}

// Init collector
func (c *Collector) Init(configPath string) {

	c.config.LoadConfig(configPath)
}

// DoDownload run download images tasks
func (c *Collector) DoDownload() {

	token := c.config.BearerToken
	imgSize := c.config.ImageSize
	retweets := c.config.IncludeRetweets
	stopDays := time.Now().AddDate(0, 0, -1*c.config.SyncLastNDays)

	for _, user := range c.config.CollectUsers {

		destPath := user.DestPath
		if destPath == "" {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				log.Fatal(err)
			}
			destPath = dir
		}

		folderName := user.FolderName
		if folderName == "" {
			folderName = user.UserID
		}

		userFolderPath := path.Join(destPath, folderName)
		if _, err := os.Stat(userFolderPath); os.IsNotExist(err) {
			os.Mkdir(userFolderPath, os.ModeDir)
		}

		userImgCnt := 0
		userDWEnd := false
		lastTweet := "0"

		for {
			tweets := getTweets(token, user.UserID, lastTweet, retweets)
			fmt.Printf("get twitter: %d\n", len(tweets))
			if len(tweets) <= 0 {
				fmt.Println("Stop task due to no tweets anymore")
				break
			} else {
				//fmt.Println(tweets)
				for _, tweet := range tweets {
					// Thu Apr 06 15:28:43 +0000 2017
					createTime, _ := time.Parse(time.RubyDate, tweet["created_at"].(string))
					if createTime.Before(stopDays) {
						fmt.Println("Stop task due to " + createTime.Format("2006-0102") + " < " + stopDays.Format("2006-0102"))
						userDWEnd = true
						break
					}

					imgURLs := extractImage(tweet)

					imgCnt := 0
					for imgURL := range imgURLs {
						fname := createTime.Format("2006-0102-150405")
						if imgCnt != 0 {
							fname += "_" + strconv.Itoa(imgCnt)
						}
						log.Println(fname)
						log.Println(imgURL)

						imgCnt++
						userImgCnt++
						if saveImage(imgURL, imgSize, path.Join(userFolderPath, fname)) {
							fmt.Printf("(%d) downloaded: %s\n", userImgCnt, fname)
						} else {
							fmt.Printf("(%d) skipped: %s\n", userImgCnt, fname)
						}
					}

					lastTweet = strconv.FormatFloat(tweet["id"].(float64), 'f', 0, 64)
				}
				if userDWEnd {
					break
				}
			}
		}
	}
}

func saveImage(imgURL, imgSize, filePath string) bool {

	destFilePath := filePath + filepath.Ext(imgURL)
	fmt.Println(destFilePath)
	if _, err := os.Stat(destFilePath); err == nil || os.IsExist(err) {
		return false
	}

	downloadPath := imgURL
	if imgSize != "" {
		downloadPath = imgURL + ":" + imgSize
	}

	resp, err := http.Get(downloadPath)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer resp.Body.Close()

	out, err := os.Create(destFilePath)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func extractImage(tweet map[string]interface{}) map[string]bool {

	imgURLs := make(map[string]bool)

	entities := tweet["entities"].(map[string]interface{})
	if media, ok := entities["media"]; ok {
		mediaList := media.([]interface{})
		for _, m := range mediaList {
			mediaObj := m.(map[string]interface{})
			if url, ok := mediaObj["media_url"]; ok {
				imgURLs[url.(string)] = true
			}
		}
	}

	if extendedEntities, ok := tweet["extended_entities"]; ok {
		entities2 := extendedEntities.(map[string]interface{})
		if media, ok := entities2["media"]; ok {
			mediaList := media.([]interface{})
			for _, m := range mediaList {
				mediaObj := m.(map[string]interface{})
				if url, ok := mediaObj["media_url"]; ok {
					imgURLs[url.(string)] = true
				}
			}
		}
	}

	return imgURLs
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

// Hello function
func Hello() string {
	return "Hello"
}
