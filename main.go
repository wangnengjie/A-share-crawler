package main

import (
	"fmt"
	"github.com/wangnengjie/A-share-spider/modules"
)

const n = 50

func main() {
	sh := modules.ReadID("./stockid/sh", "sh")
	sz := modules.ReadID("./stockid/sz", "sz")

	// update > collect > writer
	update := make(chan modules.StockMsg, 1000)
	collect := make(chan []modules.StockMsg, 10)

	go modules.WriteFile(collect)
	
}