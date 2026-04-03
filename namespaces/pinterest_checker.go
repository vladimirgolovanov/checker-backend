package namespaces

import (
	"io"
	"net/http"
	"strings"
)

type PinterestChecker struct{}

func (i *PinterestChecker) GetId() int {
	return 12
}

func (i *PinterestChecker) GetName() string {
	return "Pinterest"
}

func (i *PinterestChecker) PrepareName(name string) string {
	return name
}

func (i *PinterestChecker) ValidateName(name string) error {
	return nil
}

func (i *PinterestChecker) Check(name string, params map[string]interface{}) CheckStatus {
	url := "https://www.pinterest.com/" + name + "/"
	response, err := http.Get(url)
	if err != nil {
		return StatusFailed
	}
	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return StatusFailed
	}

	if strings.Contains(string(body), "Profile | Pinterest</title>") {
		return StatusUsed
	}

	return StatusFree
}
