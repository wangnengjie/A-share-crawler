package modules

import (
	"fmt"
	"sort"
	"time"
)

func CollectMsg(in <-chan StockMsg, out chan<- StockMsgs) {
	tick := time.NewTicker(21 * time.Second)
	var collector StockMsgs
	for {
		select {
		case msg := <-in:
			collector = append(collector, msg)
		case <-tick.C:
			n := len(collector) / 2
			sort.Sort(collector)
			fmt.Println(n)
			out <- collector[0:n]
			collector = collector[n:]
		}
	}
}
