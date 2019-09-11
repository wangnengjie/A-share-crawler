package modules

import (
	"encoding/csv"
	"fmt"
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
	f, err := os.Create("./data/" + time.Now().Format("2006-01-02_15-04-05") + ".csv")
	defer f.Close()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	writer := csv.NewWriter(f)
	defer writer.Flush()

	ticker := time.NewTicker(15 * time.Minute)
	for {
		select {
		case datas := <-in:
			writer.WriteAll(datas)
		case <-ticker.C:
			f.Close()
			writer.Flush()
			f, err = os.Create("./data/" + time.Now().Format("2006-01-02_15-04-05") + ".csv")
			writer = csv.NewWriter(f)
		}
	}
}
