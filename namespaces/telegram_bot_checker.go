package namespaces

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

type TelegramBotChecker struct {
}

func (i *TelegramBotChecker) GetId() int {
	return 10
}

func (i *TelegramBotChecker) GetName() string {
	return "Telegram bot"
}

func (i *TelegramBotChecker) PrepareName(name string) string {
	return name
}

func (i *TelegramBotChecker) ValidateName(name string) error {
	// more than 5 symbols
	if len(name) < 5 {
		return errors.New("Name must be at least 5 characters long")
	}

	// dont start with number
	if name[0] >= '0' && name[0] <= '9' {
		return errors.New("Name must not start with a number")
	}

	// only a-z, 0-9, _
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_') {
			return errors.New("Name may consist only of a-z, 0-9, and underscores")
		}
	}

	// ends with "bot"
	if !strings.HasSuffix(name, "bot") {
		return errors.New("Name must end with 'Bot'")
	}

	return nil
}

func (i *TelegramBotChecker) Check(name string, params map[string]interface{}) CheckStatus {
	url := "https://t.me/" + name + "bot"
	response, err := http.Get(url)
	if err != nil {
		return StatusFailed
	}
	defer func() { _ = response.Body.Close() }()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return StatusFailed
	}

	if strings.Contains(string(body), "<meta property=\"twitter:title\" content=\"Telegram: Contact @") {
		return StatusFree
	}

	return StatusUsed
}
