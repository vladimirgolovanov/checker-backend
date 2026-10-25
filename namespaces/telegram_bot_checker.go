package namespaces

import (
	"context"
	"errors"
	"strings"
	"time"
)

type TelegramBotChecker struct {
	fetcher Fetcher
}

func NewTelegramBotChecker(f Fetcher) *TelegramBotChecker {
	return &TelegramBotChecker{fetcher: f}
}

func (i *TelegramBotChecker) GetId() int {
	return 10
}

func (i *TelegramBotChecker) GetName() string {
	return "Telegram bot"
}

func (i *TelegramBotChecker) GetSeparateSymbols() []string {
	return []string{
		"_",
	}
}

func (i *TelegramBotChecker) PrepareName(name string) string {
	return strings.ToLower(strings.TrimSpace(name)) + "_bot"
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
			return errors.New("Name may consist only letters, numbers and underscores")
		}
	}

	// ends with "bot"
	if !strings.HasSuffix(name, "bot") {
		return errors.New("Name must end with 'bot'")
	}

	return nil
}

func (i *TelegramBotChecker) Check(ctx context.Context, name string, params map[string]interface{}) CheckStatus {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	resp, status := i.fetcher.Get(ctx, "https://t.me/"+name, nil)
	if status != 0 {
		return status
	}

	if strings.Contains(string(resp.Body), "<meta property=\"twitter:title\" content=\"Telegram: Contact @") {
		return StatusFree
	}

	return StatusUsed
}
