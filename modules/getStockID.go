package modules

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type GetStockIDMeta struct {
	ids []string
	err error
}

func GetStockID(idMap *sync.Map) ([]string, bool) {
	ch := make(chan GetStockIDMeta)
	var wg sync.WaitGroup
	wg.Add(1)
	go getSH(idMap, ch, &wg)
	wg.Add(1)
	go getSZ(idMap, ch, &wg)

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := GetStockIDMeta{ids: []string{}, err: nil}
	ok := true

	for val := range ch {
		if val.err != nil {
			LogError(val.err.Error())
			ok = false
		} else {
			result.ids = append(result.ids, val.ids...)
		}
	}

	return result.ids, ok
}

func getSH(idMap *sync.Map, out chan<- GetStockIDMeta, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get("http://www.sse.com.cn/js/common/ssesuggestdata.js")
	if err != nil {
		out <- GetStockIDMeta{ids: nil, err: err}
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		out <- GetStockIDMeta{ids: nil, err: err}
		return
	}

	reg, _ := regexp.Compile(`(60|68)[\d]{4}`)
	result := []string{}
	for _, str := range strings.Split(string(data), "\n") {
		s := reg.FindString(str)
		if len(s) == 6 {
			result = append(result, "sh"+s)
			idMap.LoadOrStore("sh"+s, "")
		}
	}
	out <- GetStockIDMeta{ids: result, err: nil}
}

func getSZ(idMap *sync.Map, out chan<- GetStockIDMeta, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get("http://www.szse.cn/api/report/ShowReport?SHOWTYPE=xlsx&CATALOGID=1110&TABKEY=tab1&random=0.46241681623081976")
	if err != nil {
		out <- GetStockIDMeta{ids: nil, err: err}
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		out <- GetStockIDMeta{ids: nil, err: err}
		return
	}

	f, err := os.Create("./data/temp.xlsx")
	defer os.Remove("./data/temp.xlsx")
	if err != nil {
		out <- GetStockIDMeta{ids: nil, err: err}
		return
	}

	f.Write(data)
	f.Close()

	xlsxFile, err := excelize.OpenFile("./data/temp.xlsx")
	if err != nil {
		out <- GetStockIDMeta{ids: nil, err: err}
		return
	}

	result := []string{}
	for _, row := range xlsxFile.GetRows("A股列表")[1:] {
		result = append(result, "sz"+row[5])
		idMap.LoadOrStore("sz"+row[5], "")
	}
	out <- GetStockIDMeta{ids: result, err: nil}
}
