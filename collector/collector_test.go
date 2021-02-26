package collector

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInit_LoadConfig_ReturnInitResult(t *testing.T) {

	tests := []struct {
		expected   bool
		configPath string
	}{
		{expected: false, configPath: "invalid config"},
		{expected: true, configPath: "valid config"},
	}

	for _, test := range tests {

		Convey(fmt.Sprintf("valid config: %t", test.expected), t, func() {
			ctl := gomock.NewController(t)
			mockConfig := NewMockIConfig(ctl)
			gomock.InOrder(
				mockConfig.EXPECT().LoadConfigFromPath(test.configPath).Return(test.expected),
			)

			c := Collector{Conf: mockConfig}
			ok := c.Init(test.configPath)
			if ok {
				So(ok, ShouldBeTrue)
			} else {
				So(ok, ShouldBeFalse)
			}
		})
	}
}

func TestDoDownload_ProcessDownload_VerifyResult(t *testing.T) {

	tests := []struct {
		expectedWithImgUrls bool
		tweetData           string
	}{
		{expectedWithImgUrls: true, tweetData: `[{"created_at":"Thu Apr 06 15:28:43 +0000 2017", "id": 100, "entities":{"media":[{"media_url":"https://i.imgur.com/IKBYy9Y.jpg"}]}, "extended_entities": {"media":[{"media_url":"https://i.imgur.com/hqGapcm.jpg"}]},"full_text":"OOOXXX"}]`},
		{expectedWithImgUrls: false, tweetData: "[]"},
	}

	for _, test := range tests {

		Convey(fmt.Sprintf("expectedWithImgUrls: %t", test.expectedWithImgUrls), t, func() {

			fakeGetTweetsReturn := make([]map[string]interface{}, 1)
			json.Unmarshal([]byte(test.tweetData), &fakeGetTweetsReturn)

			fakeGetEndTweetsReturn := make([]map[string]interface{}, 1)
			json.Unmarshal([]byte("[]"), &fakeGetEndTweetsReturn)

			ctlTwitterAPI := gomock.NewController(t)
			mockTwitterAPI := NewMockITwitterAPI(ctlTwitterAPI)

			if test.expectedWithImgUrls {
				gomock.InOrder(
					mockTwitterAPI.EXPECT().GetTweets("example token", "example user", "0", true).Return(fakeGetTweetsReturn),
					mockTwitterAPI.EXPECT().GetTweets("example token", "example user", "100", true).Return(fakeGetEndTweetsReturn),
				)
			} else {
				gomock.InOrder(
					mockTwitterAPI.EXPECT().GetTweets("example token", "example user", "0", true).Return(fakeGetTweetsReturn),
				)
			}

			var conf *Config = &Config{BearerToken: "example token", SyncLastNDays: 36500, ThreadCnt: 1, IncludeRetweets: true, CollectUsers: []UserData{{UserID: "example user", SaveDetail: true}}}

			c := Collector{Conf: conf, twitterAPI: mockTwitterAPI}
			c.Init("")
			c.DoDownload()

			if test.expectedWithImgUrls {
				So(c.userImgCnt, ShouldEqual, 2)
			} else {
				So(c.userImgCnt, ShouldEqual, 0)
			}
		})
	}
}
