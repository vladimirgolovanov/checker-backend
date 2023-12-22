package namespaces

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type TiktokChecker struct{}

func (i *TiktokChecker) GetId() int {
	return 5
}

func (i *TiktokChecker) Check(name string) CheckStatus {
	url := "https://www.tiktok.com/@" + name
	response, err := http.Get(url)
	if err != nil {
		return StatusFailed
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return StatusFailed
	}

	if strings.Contains(string(body), "uniqueId") {
		return StatusUsed
	}

	return StatusFree
}
