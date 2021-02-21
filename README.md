# twitter-images-collector

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/429545f2c9d9403cb1f89671223abc54)](https://app.codacy.com/gh/zmcx16/twitter-images-collector?utm_source=github.com&utm_medium=referral&utm_content=zmcx16/twitter-images-collector&utm_campaign=Badge_Grade_Settings)

Trace and download specified Twitter users' post images

# config.json
```
{
  "api_key": "{You twitter developer API key}",
  "api_secret": "{You twitter developer API secret}",
  "thread_cnt": 3,
  "image_size": "orig",
  "enable_log": false,
  "sync_last_n_days": 30,
  "include_retweets": true,
  "collect_users": [
	{"user_id": "{trace twitter user id}", "folder_name": "", "dest_path": ""}
  ]
}
```
Put config.json into the main program folder, set api_key, api_secret and collect_users in config.json to trace / download the specified Twitter users' post images.
* api_key & api_secret: Apply twitter developer account to get API key and API secret (https://developer.twitter.com/) 
* image_size: "large", "medium", "small", "orig", "thumb"
* sync_last_n_days: trace last N days tweet post images
* collect_users.user_id: trace tweet user ID
* collect_users.folder_name: output folder name, default is collect_users.user_id
* collect_users.dest_path: output path

# Reference
morinokami / twitter-image-downloader - (https://github.com/morinokami/twitter-image-downloader)

# License
This project is licensed under the terms of the MIT license.