package modules

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

// StockMsg : [ID,Price,Turnover,UpdateTime]
type StockMsg []string

type StockMsgs [][]string

func (v StockMsgs) Len() int      { return len(v) }
func (v StockMsgs) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v StockMsgs) Less(i, j int) bool {
	p, _ := strconv.ParseInt(v[i][3], 10, 64)
	q, _ := strconv.ParseInt(v[j][3], 10, 64)
	return p < q
}

func WriteFile(in <-chan StockMsgs) {
	name := time.Now().Format("2006-01-02_15-04-05") + ".csv"
	f, err := os.Create("./data/" + name)
	defer f.Close()
	LogInfo("create file " + name)
	if err != nil {
		LogError(err.Error())
		os.Exit(2)
	}

	writer := csv.NewWriter(f)
	defer writer.Flush()

	// 无新数据则删除创建的文件
	hasNewDate := false
	// 检测是否关闭爬虫
	getNewInFiveMinute := false
	// 检测是否收盘日，正常开盘五分钟后无新数据认为不开市，关闭爬虫
	detect := time.NewTicker(5 * time.Minute)

	for {
		select {
		case datas, ok := <-in:
			if !ok {
				writer.Flush()
				f.Close()
				detect.Stop()
				if !hasNewDate {
					os.Remove("./data/" + name)
					LogInfo("remove file " + name)
				}
				return
			}
			if !hasNewDate && len(datas) > 0 {
				hasNewDate = true
			}
			if !getNewInFiveMinute && len(datas) > 0 {
				getNewInFiveMinute = true
			}
			writer.WriteAll(datas)
		case <-detect.C:
			if getNewInFiveMinute {
				getNewInFiveMinute = false
			} else {
				SetOpenState(false)
			}
		}
	}
}
