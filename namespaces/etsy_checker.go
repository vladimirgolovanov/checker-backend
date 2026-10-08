package namespaces

import (
	"fmt"
	"unicode"
)

type EtsyChecker struct{}

func (i *EtsyChecker) GetId() int {
	return 11
}

func (i *EtsyChecker) GetName() string {
	return "Esty"
}

func (i *EtsyChecker) GetSeparateSymbols() []string {
	return []string{
		"_",
	}
}

func (i *EtsyChecker) PrepareName(name string) string {
	return name
}

// a-z, 0-9, _
// min 4, max 20
func (i *EtsyChecker) ValidateName(name string) error {
	if len(name) < 4 || len(name) > 20 {
		return fmt.Errorf("Name must be between 4 and 20 characters")
	}

	for _, c := range name {
		if !unicode.IsLower(c) && !unicode.IsDigit(c) && c != '_' {
			return fmt.Errorf("Name must contain only lowercase letters, numbers, or underscores")
		}
	}

	return nil
}

func (i *EtsyChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://www.etsy.com/shop/"+name, nil)
	if status != 0 {
		return status
	}

	if resp.StatusCode == 404 {
		return StatusFree
	}

	return StatusUsed
}
