package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/zmcx16/twitter-images-collector/collector"
)

func main() {

	fmt.Println(collector.Hello())

	logFolderPath := path.Join(filepath.Dir(os.Args[0]), "log")
	if _, err := os.Stat(logFolderPath); os.IsNotExist(err) {
		os.Mkdir(logFolderPath, os.ModeDir)
	}

	r, _ := rotatelogs.New(path.Join(logFolderPath, "%Y%m%d.log"))
	// mw := io.MultiWriter(os.Stdout, r)
	mw := io.MultiWriter(r)
	log.SetOutput(mw)
	log.SetReportCaller(true)
	log.Println(collector.Hello())
	c := collector.Collector{}
	c.Init("I:/work/WORK/GO/twitter-images-collector/config.json")
	c.DoDownload()
}
