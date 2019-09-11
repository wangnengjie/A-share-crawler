package main

import (
	"github.com/wangnengjie/A-share-spider/modules"
)

const n = 50

func main() {
	sh := modules.ReadID("./stockID/sh")
	sz := modules.ReadID("./stockID/sz")

	// update > collect > writer
	update := make(chan modules.StockMsg, 3000)
	collect := make(chan modules.StockMsgs, 3)

	for i := 0; i < len(sh); i += n {
		if i+n > len(sh) {
			go modules.Request(sh[i:], "sh", update)
		} else {
			go modules.Request(sh[i:i+n], "sh", update)
		}
	}
	for i := 0; i < len(sz); i += n {
		if i+n > len(sz) {
			go modules.Request(sz[i:], "sz", update)
		} else {
			go modules.Request(sz[i:i+n], "sz", update)
		}
	}
	go modules.CollectMsg(update, collect)

	modules.WriteFile(collect)
}
