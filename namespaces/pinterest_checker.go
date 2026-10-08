package namespaces

import (
	"fmt"
	"regexp"
	"strings"
)

// A username must contain only numbers and letters
// Usernames should be 3 to 30 characters.

type PinterestChecker struct{}

func (i *PinterestChecker) GetId() int {
	return 12
}

func (i *PinterestChecker) GetName() string {
	return "Pinterest"
}

func (i *PinterestChecker) GetSeparateSymbols() []string {
	return []string{}
}

func (i *PinterestChecker) PrepareName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// Usernames should be 3 to 30 characters.
// A username must contain only numbers and letters
// A username must have at least one letter
var pinterestNameRe = regexp.MustCompile(`^[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*$`)

func (i *PinterestChecker) ValidateName(name string) error {
	if len(name) < 3 || len(name) > 30 {
		return fmt.Errorf("Name must be between 3 and 30 characters")
	}

	if !pinterestNameRe.MatchString(name) {
		return fmt.Errorf("Name must contain only letters and numbers, and have at least one letter")
	}

	return nil
}

func (i *PinterestChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://www.pinterest.com/"+name+"/", nil)
	if status != 0 {
		return status
	}

	if strings.Contains(string(resp.Body), "Profile | Pinterest</title>") {
		return StatusUsed
	}

	return StatusFree
}
