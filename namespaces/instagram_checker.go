package namespaces

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type InstagramChecker struct{}

func (i *InstagramChecker) GetId() int {
	return 0
}

func (i *InstagramChecker) GetName() string {
	return "Instagram"
}

func (i *InstagramChecker) PrepareName(name string) string {
	return name
}

func (i *InstagramChecker) ValidateName(name string) error {
	return nil
}

func (i *InstagramChecker) Check(name string, params map[string]interface{}) CheckStatus {
	url := "https://www.instagram.com/" + name + "/"
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

	// req.Header.Set("Host", "www.instagram.com")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:125.0) Gecko/20100101 Firefox/125.0")
	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	// req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	// req.Header.Set("DNT", "1")
	// req.Header.Set("Sec-GPC", "1")
	// req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("Upgrade-Insecure-Requests", "1")
	// req.Header.Set("Sec-Fetch-Dest", "document")
	// req.Header.Set("Sec-Fetch-Mode", "navigate")
	// req.Header.Set("Sec-Fetch-Site", "none")
	// req.Header.Set("Sec-Fetch-User", "?1")

	// req.Header.Set("Host", "www.instagram.com")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:125.0) Gecko/20100101 Firefox/125.0")
	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	// req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	// req.Header.Set("DNT", "1")
	// req.Header.Set("Alt-Used", "www.instagram.com")
	// req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("Cookie", "csrftoken=x2aSCAJsI5vLLIXq2smnzT")
	// req.Header.Set("Upgrade-Insecure-Requests", "1")
	// req.Header.Set("Sec-Fetch-Dest", "document")
	// req.Header.Set("Sec-Fetch-Mode", "navigate")
	// req.Header.Set("Sec-Fetch-Site", "same-origin")
	// req.Header.Set("TE", "trailers")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("failed on do request: ", err)
		return StatusFailed
	}
	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("failed on read response body: ", err)
		return StatusFailed
	}

	if strings.Contains(string(body), "{\"username\":\""+name+"\"}") {
		fmt.Println("used")
		return StatusUsed
	}

	if !strings.Contains(string(body), "\"url\":\"\\/"+name+"\\/\"") {
		fmt.Println("failed on finding url string")
		return StatusFailed
	}

	fmt.Println("free")
	return StatusFree

	/*
		url := "http://localhost:3000/api/v1/check-username"
		data := map[string]interface{}{
			"username": name,
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Ошибка при кодировании JSON данных:", err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			// временно забьем на ошибку
			fmt.Println("Ошибка при выполнении запроса:", err)
		}
		req.Header.Set("Authorization", "your_auth_token_here")
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка при отправке HTTP запроса:", err)
		}
		defer resp.Body.Close()

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			fmt.Println("Ошибка при чтении JSON ответа:", err)
		}

		return response["success"] == true
	*/
}
