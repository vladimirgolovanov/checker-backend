package namespaces

import (
	"net/http"
)

type GithubChecker struct{}

func (i *GithubChecker) GetId() int {
	return 8
}

func (i *GithubChecker) Check(name string) CheckStatus {
	url := "https://github.com/" + name
	response, err := http.Get(url)
	if err != nil {
		return StatusFailed
		// todo: записать ошибку в лог
	}
	defer response.Body.Close()

	if response.StatusCode != 404 {
		return StatusUsed
	}

	return StatusFree
}
