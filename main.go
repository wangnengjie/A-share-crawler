package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wangnengjie/A-share-crawler/modules"
)

func main() {
	v := flag.Int("v", 50, "单次请求股票数，设置越低并发数越高")
	t := flag.Int("t", 15, "单文件记录时间区间，单位分钟")
	flag.Parse()

	if *v < 1 {
		fmt.Fprintln(os.Stderr, "参数v应大于0")
		os.Exit(1)
	}
	if *t < 1 {
		fmt.Fprintln(os.Stderr, "参数t应大于0")
		os.Exit(1)
	}

	sh := modules.ReadID("./stockID/sh")
	sz := modules.ReadID("./stockID/sz")

	// update > collect > writer
	update := make(chan modules.StockMsg, 3000)
	collect := make(chan modules.StockMsgs, 3)

	for i := 0; i < len(sh); i += *v {
		if i+*v > len(sh) {
			go modules.Request(sh[i:], "sh", update)
		} else {
			go modules.Request(sh[i:i+*v], "sh", update)
		}
	}
	for i := 0; i < len(sz); i += *v {
		if i+*v > len(sz) {
			go modules.Request(sz[i:], "sz", update)
		} else {
			go modules.Request(sz[i:i+*v], "sz", update)
		}
	}
	go modules.CollectMsg(update, collect)

	modules.WriteFile(collect, t)
}
