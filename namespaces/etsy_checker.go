package namespaces

import (
	"fmt"
	"net/http"
)

type EstyChecker struct {
}

func (i *EstyChecker) GetId() int {
	return 11
}

func (i *EstyChecker) GetName() string {
	return "Esty"
}

func (i *EstyChecker) PrepareName(name string) string {
	return name
}

// если содержит a-z, 0-9, _ и длиной от 4 до 20 символов
// если нет, то возвращаем ошибку
func (i *EstyChecker) ValidateName(name string) error {
	return nil
}

func (i *EstyChecker) Check(name string, params map[string]interface{}) CheckStatus {
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
	defer response.Body.Close()

	if response.StatusCode == 404 {
		return StatusFree
	}

	return StatusUsed
}
