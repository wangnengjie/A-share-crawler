package modules

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Request(ids []string, idMap *sync.Map, out chan<- StockMsg) {
	url := ""
	for _, id := range ids {
		url += id + ","
	}
	for ; ; time.Sleep(1 * time.Second) {
		if !GetOpenState() {
			return
		}

		resp, err := http.Get("http://hq.sinajs.cn/list=" + url)

		if err != nil {
			LogError(err.Error())
			continue
		}
		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			LogError(err.Error())
			continue
		}

		msgs := strings.Split(string(data), "\n")

		for index, msg := range msgs[:len(msgs)-1] {
			go func(index int, s string) {
				pre, _ := idMap.Load(ids[index])
				if pre != nil && s != pre {
					idMap.Store(ids[index], s)
					m := getStockMsg(s, ids[index])
					if m[1] == "0.000" {
						return
					}
					if !GetOpenState() {
						return
					}
					out <- m
				}
			}(index, msg)
		}
	}
}

func getStockMsg(s string, id string) StockMsg {
	slice := strings.Split(s, ",")
	return StockMsg{
		id[2:],
		slice[3],
		slice[8],
		strconv.FormatInt(getTime(&slice[30], &slice[31]), 10),
	}
}

func getTime(date *string, t *string) int64 {
	t_p, err := time.Parse("2006-01-02 15:04:05", *date+" "+*t)
	if err != nil {
		LogError(err.Error())
	}
	return t_p.Unix()
}
