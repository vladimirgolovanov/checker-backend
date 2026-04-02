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

func (i *TiktokChecker) GetName() string {
	return "Tiktok"
}

func (i *TiktokChecker) PrepareName(name string) string {
	return name
}

func (i *TiktokChecker) ValidateName(name string) error {
	return nil
}

func (i *TiktokChecker) Check(name string, params map[string]interface{}) CheckStatus {
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
