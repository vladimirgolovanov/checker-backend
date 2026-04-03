package namespaces

import (
	"fmt"
	"net/http"
)

type EtsyChecker struct {
}

func (i *EtsyChecker) GetId() int {
	return 11
}

func (i *EtsyChecker) GetName() string {
	return "Esty"
}

func (i *EtsyChecker) PrepareName(name string) string {
	return name
}

// если содержит a-z, 0-9, _ и длиной от 4 до 20 символов
// если нет, то возвращаем ошибку
func (i *EtsyChecker) ValidateName(name string) error {
	return nil
}

func (i *EtsyChecker) Check(name string, params map[string]interface{}) CheckStatus {
	url := "https://www.etsy.com/shop/" + name
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("failed on create request: ", err)
		return StatusFailed
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed on do request: ", err)
		return StatusFailed
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode == 404 {
		return StatusFree
	}

	return StatusUsed
}
