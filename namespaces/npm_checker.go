package namespaces

import (
	"io"
	"net/http"
)

type NpmChecker struct {
}

func (i *NpmChecker) GetId() int {
	return 7
}

func (i *NpmChecker) GetName() string {
	return "Npm"
}

func (i *NpmChecker) PrepareName(name string) string {
	return name
}

func (i *NpmChecker) ValidateName(name string) error {
	return nil
}

func (i *NpmChecker) Check(name string, params map[string]interface{}) CheckStatus {
	url := "https://www.npmjs.com/package/" + name
	response, err := http.Get(url)
	if err != nil {
		return StatusFailed
	}
	defer func() {
		_, _ = io.Copy(io.Discard, response.Body)
		_ = response.Body.Close()
	}()

	if response.StatusCode == 404 {
		return StatusFree
	}

	return StatusUsed
}
