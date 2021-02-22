package collector

import (
	"container/list"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Collector struct
type Collector struct {
	config           Config
	userImgCnt       int
	muxThreadStopped sync.RWMutex
	muxUserImgCnt    sync.RWMutex
	muxLastTwitterID sync.RWMutex
	muxTweet         sync.RWMutex
}

// Init collector
func (c *Collector) Init(configPath string) bool {

	if !c.config.LoadConfig(configPath) {
		fmt.Printf("Load Config file failed (%s)\n", configPath)
		return false
	}
	return true
}

// DoDownload run download images tasks
func (c *Collector) DoDownload() {

	token := c.config.BearerToken
	threadCnt := c.config.ThreadCnt
	imgSize := c.config.ImageSize
	retweets := c.config.IncludeRetweets
	stopDays := time.Now().AddDate(0, 0, -1*c.config.SyncLastNDays)

	for _, user := range c.config.CollectUsers {

		destPath := string(user.DestPath)
		destPath = strings.Trim(destPath, "\"") // remove json.RawMessage start and end ""
		if string(destPath) == "" {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				log.Error(err)
			}
			destPath = dir
		}

		folderName := string(user.FolderName)
		folderName = strings.Trim(folderName, "\"") // remove json.RawMessage start and end ""
		if folderName == "" {
			folderName = user.UserID
		}

		userFolderPath := path.Join(destPath, folderName)
		if _, err := os.Stat(userFolderPath); os.IsNotExist(err) {
			os.MkdirAll(userFolderPath, os.ModeDir)
		}

		userDWEnd := false
		lastTweet := "0"

		for {
			twitterAPI := &TwitterAPI{Client: &http.Client{}}
			tweets := twitterAPI.GetTweets(token, user.UserID, lastTweet, retweets)
			fmt.Printf("get twitter list: %d, (%s)\n", len(tweets), lastTweet)
			if len(tweets) <= 0 {
				fmt.Printf("Stop task due to no tweets anymore (%s)\n", user.UserID)
				break
			} else {
				//fmt.Println(tweets)
				lastTweet2 := lastTweet
				lastTweet, userDWEnd = c.dwTweetImgs(tweets, stopDays, threadCnt, user.SaveDetail, imgSize, user.UserID, userFolderPath)
				if userDWEnd || lastTweet == lastTweet2 {
					fmt.Println("Stop download task")
					break
				}
			}
		}
	}
}

func (c *Collector) dwTweetImgs(tweets []map[string]interface{}, stopDays time.Time, threadCnt int, saveDetail bool, imgSize, userID, userFolderPath string) (lastTweet string, userDWEnd bool) {

	lastTweetFloat := math.MaxFloat64
	var threadTerminatedList = make([]bool, threadCnt)
	tweetIndexQueue := list.New()
	for i := range tweets {
		tweetIndexQueue.PushBack(i)
	}

	var wg sync.WaitGroup

	for ti := 0; ti < threadCnt; ti++ {

		wg.Add(1)
		go func(ti int) {
			defer wg.Done()

			for {
				c.muxTweet.Lock()

				if tweetIndexQueue.Len() == 0 {
					fmt.Printf("[T%d] thread end\n", ti)
					c.muxTweet.Unlock()
					return
				}

				i := tweetIndexQueue.Front()
				tweetIndexQueue.Remove(i)
				tweet := tweets[i.Value.(int)]

				c.muxTweet.Unlock()

				terminated := false
				// Thu Apr 06 15:28:43 +0000 2017
				createTime, _ := time.Parse(time.RubyDate, tweet["created_at"].(string))
				if createTime.Before(stopDays) {
					fmt.Printf("[T%d] Stop task due to %s < %s\n", ti, createTime.Format("2006-0102"), stopDays.Format("2006-0102"))
					terminated = true
					userDWEnd = true
				}

				if terminated {
					c.muxThreadStopped.Lock()
					threadTerminatedList[ti] = true
					c.muxThreadStopped.Unlock()
					return
				}

				imgURLs, tweetText := extractImage(tweet, saveDetail)

				imgCnt := 0
				fname := createTime.Format("2006-0102-150405")
				for imgURL := range imgURLs {
					imgFname := fname
					if imgCnt != 0 {
						imgFname += "_" + strconv.Itoa(imgCnt)
					}
					log.Println(imgFname)
					log.Println(imgURL)

					imgCnt++
					if saveImage(imgURL, imgSize, path.Join(userFolderPath, imgFname)) {
						c.muxUserImgCnt.Lock()
						c.userImgCnt++
						fmt.Printf("(%d) downloaded image: %s\n", c.userImgCnt, imgFname)
						c.muxUserImgCnt.Unlock()
					} else {
						c.muxUserImgCnt.Lock()
						c.userImgCnt++
						fmt.Printf("(%d) skipped image: %s\n", c.userImgCnt, imgFname)
						c.muxUserImgCnt.Unlock()
					}
				}

				if imgCnt > 0 && saveDetail {
					if saveDetailInfo(tweetText, path.Join(userFolderPath, fname)) {
						fmt.Printf("downloaded detail: %s\n", fname)
					} else {
						fmt.Printf("skipped detail: %s\n", fname)
					}
				}

				c.muxLastTwitterID.Lock()
				lastTweetFloat = math.Min(lastTweetFloat, tweet["id"].(float64))
				c.muxLastTwitterID.Unlock()
			}
		}(ti)
	}
	wg.Wait()

	threadAllTerminated := true
	for ti := 0; ti < threadCnt; ti++ {
		threadAllTerminated = threadAllTerminated && threadTerminatedList[ti]
	}

	if threadAllTerminated {
		fmt.Println("all thread terminated")
		userDWEnd = true
	}

	return strconv.FormatFloat(lastTweetFloat, 'f', 0, 64), userDWEnd
}

func saveImage(imgURL, imgSize, filePath string) bool {

	destImgFilePath := filePath + filepath.Ext(imgURL)
	if _, err := os.Stat(destImgFilePath); err == nil || os.IsExist(err) {
		return false
	}

	downloadPath := imgURL
	if imgSize != "" {
		downloadPath = imgURL + ":" + imgSize
	}

	resp, err := http.Get(downloadPath)
	if err != nil {
		log.Error(err)
		return false
	}
	defer resp.Body.Close()

	out, err := os.Create(destImgFilePath)
	if err != nil {
		log.Error(err)
		return false
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Error(err)
		return false
	}

	return true
}

func saveDetailInfo(tweetText, filePath string) bool {

	destDetailFilePath := filePath + ".txt"
	if _, err := os.Stat(destDetailFilePath); err == nil || os.IsExist(err) {
		return false
	}

	err := ioutil.WriteFile(destDetailFilePath, []byte(tweetText), 0644)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func extractImage(tweet map[string]interface{}, saveDetail bool) (imgURLs map[string]bool, tweetText string) {

	imgURLs = make(map[string]bool)

	if saveDetail {
		tweetText = tweet["full_text"].(string)
	}

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

	return
}
