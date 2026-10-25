package namespaces

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode"
)

// Usernames can't contain numbers only
// Only lowercase letters, numbers, underscores, or periods are allowed
// Usernames can't end with a period
// min 2, max 24

type TiktokChecker struct {
	fetcher Fetcher
}

func NewTiktokChecker(f Fetcher) *TiktokChecker {
	return &TiktokChecker{fetcher: f}
}

func (i *TiktokChecker) GetId() int {
	return 5
}

func (i *TiktokChecker) GetName() string {
	return "Tiktok"
}

func (i *TiktokChecker) GetSeparateSymbols() []string {
	return []string{
		"_",
		".",
	}
}

func (i *TiktokChecker) PrepareName(name string) string {
	return name
}

func (i *TiktokChecker) ValidateName(name string) error {
	if len(name) < 2 || len(name) > 24 {
		return fmt.Errorf("Name must be between 2 and 24 characters")
	}

	for _, c := range name {
		if !unicode.IsLower(c) && !unicode.IsDigit(c) && c != '_' && c != '.' {
			return fmt.Errorf("Name must contain only lowercase letters, numbers, underscores, or periods")
		}
	}

	if name[len(name)-1] == '.' {
		return fmt.Errorf("Name must not end with a period")
	}

	return nil
}

func (i *TiktokChecker) Check(ctx context.Context, name string, params map[string]interface{}) CheckStatus {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	resp, status := i.fetcher.Get(ctx, "https://www.tiktok.com/@"+name, nil)
	if status != 0 {
		return status
	}

	if strings.Contains(string(resp.Body), "uniqueId") {
		return StatusUsed
	}

	return StatusFree
}
