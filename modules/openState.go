package modules

import (
	"fmt"
	"sync"
	"time"
)

var open bool
var mu sync.RWMutex

func GetOpenState() bool {
	mu.RLock()
	defer mu.RUnlock()
	return open
}

func SetOpenState(state bool) {
	LogInfo(fmt.Sprintln("set open state", state))
	mu.Lock()
	defer mu.Unlock()
	open = state
}

func InitializeTimer(t5, t6, t7, t8 **time.Timer) {
	now := time.Now()
	// 9:28
	t1 := time.Date(now.Year(), now.Month(), now.Day(), 9, 28, 0, 0, now.Location())
	// 11:32
	t2 := t1.Add(2*time.Hour + 4*time.Minute)
	// 12:58
	t3 := t2.Add(1*time.Hour + 26*time.Minute)
	// 15:32 大宗交易在15：00 - 15：30
	t4 := t3.Add(2*time.Hour + 34*time.Minute)

	switch {
	case now.Before(t1):
		open = false
		*t5 = time.NewTimer(t1.Sub(now))
		*t6 = time.NewTimer(t2.Sub(now))
		*t7 = time.NewTimer(t3.Sub(now))
		*t8 = time.NewTimer(t4.Sub(now))
	case now.After(t1) && now.Before(t2):
		open = true
		*t5 = time.NewTimer(24*time.Hour - now.Sub(t1))
		*t6 = time.NewTimer(t2.Sub(now))
		*t7 = time.NewTimer(t3.Sub(now))
		*t8 = time.NewTimer(t4.Sub(now))
	case now.After(t2) && now.Before(t3):
		open = false
		*t5 = time.NewTimer(24*time.Hour - now.Sub(t1))
		*t6 = time.NewTimer(24*time.Hour - now.Sub(t2))
		*t7 = time.NewTimer(t3.Sub(now))
		*t8 = time.NewTimer(t4.Sub(now))
	case now.After(t3) && now.Before(t4):
		open = true
		*t5 = time.NewTimer(24*time.Hour - now.Sub(t1))
		*t6 = time.NewTimer(24*time.Hour - now.Sub(t2))
		*t7 = time.NewTimer(24*time.Hour - now.Sub(t3))
		*t8 = time.NewTimer(t4.Sub(now))
	case now.After(t4):
		*t5 = time.NewTimer(24*time.Hour - now.Sub(t1))
		*t6 = time.NewTimer(24*time.Hour - now.Sub(t2))
		*t7 = time.NewTimer(24*time.Hour - now.Sub(t3))
		*t8 = time.NewTimer(24*time.Hour - now.Sub(t4))
	}
}
