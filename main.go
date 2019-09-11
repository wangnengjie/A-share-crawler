package main

import (
	"github.com/wangnengjie/A-share-spider/modules"
)

func main() {
	sh := modules.ReadID("./stockid/sh")
	sz := modules.ReadID("./stockid/sz")

	// update > collect > writer
	update := make(chan modules.StockMsg, 50)
	collect := make(chan modules.StockMsgs, 10)

	for _, id := range sh {
		go modules.Request(id, "sh", update)
	}
	for _, id := range sz {
		go modules.Request(id, "sz", update)
	}
	go modules.CollectMsg(update, collect)

	modules.WriteFile(collect)
}
