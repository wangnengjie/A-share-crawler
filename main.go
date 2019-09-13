package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/wangnengjie/A-share-crawler/modules"
)

func init() {
	go modules.StartLog()
	modules.CreateDataDir()
}

func main() {
	v := flag.Int("v", 50, "单次请求股票数，设置越低并发数越高")
	t := flag.Int("t", 15, "单文件记录时间区间，单位分钟")
	flag.Parse()
	testFlag(v, t)

	// 四个开盘时间
	var t1, t2, t3, t4 *time.Timer

	modules.InitializeTimer(t1, t2, t3, t4)

	var idMap sync.Map
	ids, ok := modules.GetStockID(&idMap)
	if !ok {
		modules.LogError("获取股票列表失败，请尝试重写获取爬虫")
	}

	if modules.GetOpenState() {
		start(&ids, &idMap, v, t, false)
	}

	for {
		select {
		case <-t1.C:
			modules.SetOpenState(true)
			start(&ids, &idMap, v, t, true)
			t1.Reset(24 * time.Hour)
		case <-t2.C:
			modules.SetOpenState(false)
			t2.Reset(24 * time.Hour)
		case <-t3.C:
			modules.SetOpenState(true)
			start(&ids, &idMap, v, t, true)
			t3.Reset(24 * time.Hour)
		case <-t4.C:
			modules.SetOpenState(false)
			t4.Reset(24 * time.Hour)
		}
	}

	// sh := modules.ReadID("./stockID/sh")
	// sz := modules.ReadID("./stockID/sz")

	// // update > collect > writer
	// update := make(chan modules.StockMsg, 3000)
	// collect := make(chan modules.StockMsgs, 3)

	// for i := 0; i < len(sh); i += *v {
	// 	if i+*v > len(sh) {
	// 		go modules.Request(sh[i:], "sh", update)
	// 	} else {
	// 		go modules.Request(sh[i:i+*v], "sh", update)
	// 	}
	// }
	// for i := 0; i < len(sz); i += *v {
	// 	if i+*v > len(sz) {
	// 		go modules.Request(sz[i:], "sz", update)
	// 	} else {
	// 		go modules.Request(sz[i:i+*v], "sz", update)
	// 	}
	// }
	// go modules.CollectMsg(update, collect)

	// modules.WriteFile(collect, t)
}

func start(ids *[]string, idMap *sync.Map, v *int, t *int, reget bool) {
	// 需要重新获取列表
	if reget {
		_ids, ok := modules.GetStockID(idMap)
		if ok {
			*ids = _ids
		}
	}
	// update > collect > writer
	update := make(chan modules.StockMsg, 3000)
	collect := make(chan modules.StockMsgs, 3)
	for i := 0; i < len(*ids); i += *v {
		if i+*v > len(*ids) {
			go modules.Request((*ids)[i:], idMap, update)
		} else {
			go modules.Request((*ids)[i:i+*v], idMap, update)
		}
	}
	go modules.CollectMsg(update, collect)
	modules.WriteFile(collect, t)
}

func testFlag(v *int, t *int) {
	if *v < 1 {
		fmt.Fprintln(os.Stderr, "参数v应大于0")
		os.Exit(1)
	}
	if *v > 900 {
		fmt.Fprintln(os.Stderr, "参数v设置过大")
		os.Exit(1)
	}
	if *t < 1 {
		fmt.Fprintln(os.Stderr, "参数t应大于0")
		os.Exit(1)
	}
}
