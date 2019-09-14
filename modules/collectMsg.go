package modules

import (
	"sort"
	"time"
)

func CollectMsg(in <-chan StockMsg, out chan<- StockMsgs) {
	tick := time.NewTicker(30 * time.Second)
	collector := StockMsgs{}
	for {
		if !GetOpenState() {
			for len(in) != 0 {
				collector = append(collector, <-in)
			}
			sort.Sort(collector)
			out <- collector
			close(out)
			return
		}
		select {
		case msg := <-in:
			collector = append(collector, msg)
		case <-tick.C:
			n := len(collector) / 2
			sort.Sort(collector)
			out <- collector[:n]
			collector = collector[n:]
		}
	}
}
