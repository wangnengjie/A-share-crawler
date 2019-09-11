package modules

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// ReadID : read stock id from sh sz
func ReadID(path string) []string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("fail to open file %s : %s", path, err)
	}
	ids := strings.Split(string(data), ",")
	return ids[:len(ids)-1]
}
