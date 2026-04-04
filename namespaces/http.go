package namespaces

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

type HttpResponse struct {
	StatusCode int
	Body       []byte
}

// Get выполняет GET-запрос с таймаутом и экранированием пути.
// Возвращает StatusPending при таймауте, StatusFailed при других ошибках.
func Get(rawURL string, headers map[string]string) (*HttpResponse, CheckStatus) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, StatusFailed
	}
	parsedURL.Path = parsedURL.Path
	parsedURL.RawPath = parsedURL.EscapedPath()

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, StatusFailed
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	response, err := httpClient.Do(req)
	if err != nil {
		if isTimeout(err) {
			return nil, StatusPending
		}
		return nil, StatusFailed
	}
	defer func() {
		_, _ = io.Copy(io.Discard, response.Body)
		_ = response.Body.Close()
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, StatusFailed
	}

	return &HttpResponse{
		StatusCode: response.StatusCode,
		Body:       body,
	}, 0
}

func isTimeout(err error) bool {
	type timeout interface {
		Timeout() bool
	}
	if t, ok := err.(timeout); ok {
		return t.Timeout()
	}
	return false
}
