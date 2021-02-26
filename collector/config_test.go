package collector

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadConfigFromReader_GenBearerTokenFailed_ReturnFalse(t *testing.T) {
	Convey("GenBearerTokenFailed", t, func() {
		var buffer bytes.Buffer
		buffer.WriteString("{}")
		var c *Config = &Config{}
		ok := c.LoadConfigFromReader(&buffer)

		So(ok, ShouldBeFalse)
	})
}

func TestLoadConfigFromReader_GetValidConfig_ReturnTrue(t *testing.T) {
	Convey("GetValidConfig", t, func() {
		var buffer bytes.Buffer
		buffer.WriteString(`{
			"api_key": "OOO",
			"api_secret": "XXX",
			"thread_cnt": 1,
			"image_size": "orig",
			"enable_log": false,
			"sync_last_n_days": 365,
			"include_retweets": true,
			"collect_users": [
				{"user_id": "oekakisurude12", "save_detail": true, "folder_name": "oekakisurude12", "dest_path": "K:/梗圖"}
			]
		}`)

		ctl := gomock.NewController(t)
		mockTwitterAPI := NewMockITwitterAPI(ctl)
		gomock.InOrder(
			mockTwitterAPI.EXPECT().GenBearerToken("OOO", "XXX").Return("OOOXXX"),
		)

		var c *Config = &Config{twitterAPI: mockTwitterAPI}
		ok := c.LoadConfigFromReader(&buffer)

		So(ok, ShouldBeTrue)
		So(c.BearerToken, ShouldEqual, "OOOXXX")
	})
}
