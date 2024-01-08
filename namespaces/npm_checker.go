package namespaces

import (
	"net/http"
)

type NpmChecker struct {
}

func (i *NpmChecker) GetId() int {
	return 7
}

func (i *NpmChecker) Check(name string) CheckStatus {
	url := "https://www.npmjs.com/package/" + name
	response, err := http.Get(url)
	if err != nil {
		return StatusFailed
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		return StatusFree
	}

	return StatusUsed
}
