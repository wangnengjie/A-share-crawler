package modules

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Request(id string, prefix string, out chan<- StockMsg) {
	pre := ""
	for ; ; time.Sleep(1 * time.Second) {
		resp, err := http.Get("http://hq.sinajs.cn/list=" + prefix + id)

		if err != nil {
			fmt.Fprintf(os.Stderr, "request error %s", err)
			continue
		}
		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s", err)
			continue
		}

		msg := string(data)
		if msg != pre {
			pre = msg
			out <- getStockMsg(msg, id)
		}
	}
}

func getStockMsg(s string, id string) StockMsg {
	slice := strings.Split(s, ",")
	return StockMsg{
		id,
		slice[3],
		slice[8],
		strconv.FormatInt(getTime(&slice[30], &slice[31])/1000, 10),
	}
}

func getTime(date *string, t *string) int64 {
	t_p, err := time.Parse("2006-01-02 15:04:05", *date+" "+*t)
	if err != nil {
		fmt.Printf("time parse err %s", err)
	}
	return t_p.Unix()
}
