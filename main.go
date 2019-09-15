package main

import (
	"flag"
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
	flag.Parse()
	testFlag(v)

	// 四个开盘时间
	var t1, t2, t3, t4 *time.Timer

	modules.InitializeTimer(&t1, &t2, &t3, &t4)

	var idMap sync.Map
	ids, ok := modules.GetStockID(&idMap)
	if !ok {
		modules.LogError("获取股票列表失败，请尝试重写获取爬虫")
	}

	if modules.GetOpenState() {
		start(&ids, &idMap, v, false)
	}

	for {
		select {
		case <-t1.C:
			modules.SetOpenState(true)
			start(&ids, &idMap, v, true)
			t1.Reset(24 * time.Hour)
		case <-t2.C:
			modules.SetOpenState(false)
			t2.Reset(24 * time.Hour)
		case <-t3.C:
			modules.SetOpenState(true)
			start(&ids, &idMap, v, true)
			t3.Reset(24 * time.Hour)
		case <-t4.C:
			modules.SetOpenState(false)
			t4.Reset(24 * time.Hour)
		}
	}
}

func start(ids *[]string, idMap *sync.Map, v *int, reget bool) {
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
	// todo: 是否需要多个collector
	go modules.CollectMsg(update, collect)
	go modules.WriteFile(collect)
}

func testFlag(v *int) {
	if *v < 1 {
		modules.LogError("参数v应大于0")
		os.Exit(1)
	}
	if *v > 900 {
		modules.LogError("参数v设置过大")
		os.Exit(1)
	}
}
