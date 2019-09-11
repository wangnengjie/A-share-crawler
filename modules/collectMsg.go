package modules

import (
	"sort"
	"time"
)

func CollectMsg(in <-chan StockMsg, out chan<- StockMsgs) {
	tick := time.NewTicker(30 * time.Second)
	collector := StockMsgs{}
	for {
		select {
		case msg := <-in:
			collector = append(collector, msg)
		case <-tick.C:
			sort.Sort(collector)
			out <- collector
			collector = StockMsgs{}
		}
	}
}
