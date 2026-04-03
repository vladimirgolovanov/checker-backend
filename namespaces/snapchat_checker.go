package namespaces

import (
	"io"
	"net/http"
	"strings"
)

type SnapchatChecker struct{}

func (i *SnapchatChecker) GetId() int {
	return 6
}

func (i *SnapchatChecker) GetName() string {
	return "Snapchat"
}

func (i *SnapchatChecker) PrepareName(name string) string {
	return name
}

func (i *SnapchatChecker) ValidateName(name string) error {
	return nil
}

func (i *SnapchatChecker) Check(name string, params map[string]interface{}) CheckStatus {
	url := "https://www.snapchat.com/add/" + name
	response, err := http.Get(url)
	if err != nil {
		return StatusFailed
	}
	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return StatusFailed
	}

	if strings.Contains(string(body), "<title data-react-helmet=\"true\">Snapchat</title>") {
		return StatusFree
	}

	return StatusUsed
}
