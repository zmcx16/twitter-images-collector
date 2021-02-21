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

	fmt.Println("*** twitter-images-collector start ***")

	procDirPath := filepath.Dir(os.Args[0])
	logFolderPath := path.Join(procDirPath, "log")
	if _, err := os.Stat(logFolderPath); os.IsNotExist(err) {
		os.Mkdir(logFolderPath, os.ModeDir)
	}

	r, _ := rotatelogs.New(path.Join(logFolderPath, "%Y%m%d.log"))
	// mw := io.MultiWriter(os.Stdout, r)
	mw := io.MultiWriter(r)
	log.SetOutput(mw)
	log.SetReportCaller(true)
	c := collector.Collector{}

	if !c.Init(path.Join(procDirPath, "config.json")) {
		fmt.Println("twitter-images-collector Init failed")
		os.Exit(-1)
	}

	c.DoDownload()

	fmt.Println("*** twitter-images-collector end ***")
}
