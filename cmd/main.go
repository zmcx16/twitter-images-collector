package main

import (
	"fmt"

	"github.com/zmcx16/twitter-images-collector/collector"
)

func main() {

	fmt.Println(collector.Hello())

	c := collector.Collector{}
	c.Init("I:/work/WORK/GO/twitter-images-collector/config.json")
}
