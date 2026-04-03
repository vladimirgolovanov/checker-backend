package namespaces

import (
	"net/http"
)

type GithubChecker struct{}

func (i *GithubChecker) GetId() int {
	return 8
}

func (i *GithubChecker) GetName() string {
	return "Github"
}

func (i *GithubChecker) PrepareName(name string) string {
	return name
}

func (i *GithubChecker) ValidateName(name string) error {
	return nil
}

func (i *GithubChecker) Check(name string, params map[string]interface{}) CheckStatus {
	url := "https://github.com/" + name
	response, err := http.Get(url)
	if err != nil {
		return StatusFailed
		// todo: записать ошибку в лог
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != 404 {
		return StatusUsed
	}

	return StatusFree
}
