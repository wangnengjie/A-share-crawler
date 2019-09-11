package modules

import (
	"github.com/360EntSecGroup-Skylar/excelize"
)

// StockMsg : stock update
type StockMsg struct {
	ID         string
	Price      string
	Turnover   string
	UpdateTime int
}

// WriteFile : write data to xlsx
func WriteFile(in <-chan []StockMsg) {
	f := excelize.NewFile()
	var count int64 = 1
	for datas := range in {
		for _, data := range datas {
			f.SetSheetRow("sheet1", "A"+string(count), &[]interface{}{
				data.ID,
				data.Price,
				data.Turnover,
				data.UpdateTime,
			})
			count++
		}
	}
}
