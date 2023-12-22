package namespaces

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type SnapchatChecker struct{}

func (i *SnapchatChecker) GetId() int {
	return 6
}

func (i *SnapchatChecker) Check(name string) CheckStatus {
	url := "https://www.snapchat.com/add/" + name
	response, err := http.Get(url)
	if err != nil {
		return StatusFailed
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return StatusFailed
	}

	if strings.Contains(string(body), "<title data-react-helmet=\"true\">Snapchat</title>") {
		return StatusFree
	}

	return StatusUsed
}
