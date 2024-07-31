package namespaces

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type PinterestChecker struct{}

func (i *PinterestChecker) GetId() int {
	return 12
}

func (i *PinterestChecker) Check(name string) CheckStatus {
	url := "https://www.pinterest.com/" + name + "/"
	response, err := http.Get(url)
	if err != nil {
		return StatusFailed
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return StatusFailed
	}

	if strings.Contains(string(body), "Profile | Pinterest</title>") {
		return StatusUsed
	}

	return StatusFree
}
