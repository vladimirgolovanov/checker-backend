package namespaces

import (
	"errors"
	"strings"
)

type TelegramChecker struct{}

func (i *TelegramChecker) GetId() int {
	return 9
}

func (i *TelegramChecker) GetName() string {
	return "Telegram"
}

func (i *TelegramChecker) GetSeparateSymbols() []string {
	return []string{
		"_",
	}
}

func (i *TelegramChecker) PrepareName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// Usernames are case-insensitive, must be at least 5-characters long, and may consist only of a-z, 0–9, and underscores.
func (i *TelegramChecker) ValidateName(name string) error {
	if len(name) < 5 {
		return errors.New("Name must be at least 5 characters long")
	}

	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_') {
			return errors.New("Name may consist only letters, numbers and underscores")
		}
	}

	return nil
}

func (i *TelegramChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://t.me/"+name, nil)
	if status != 0 {
		return status
	}

	if strings.Contains(string(resp.Body), "<meta property=\"twitter:title\" content=\"Telegram: Contact @") {
		return StatusFree
	}

	return StatusUsed
}
