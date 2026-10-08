package namespaces

import (
	"fmt"
	"regexp"
	"strings"
)

type SnapchatChecker struct{}

func (i *SnapchatChecker) GetId() int {
	return 6
}

func (i *SnapchatChecker) GetName() string {
	return "Snapchat"
}

func (i *SnapchatChecker) GetSeparateSymbols() []string {
	return []string{
		"-",
		"_",
		".",
	}
}

func (i *SnapchatChecker) PrepareName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))

	// only one of - _ . in username
	foundSpecial := false
	result := strings.Builder{}
	for _, c := range name {
		if c == '-' || c == '_' || c == '.' {
			if !foundSpecial {
				foundSpecial = true
				result.WriteRune(c)
			}
		} else {
			result.WriteRune(c)
		}
	}
	return result.String()
}

// Usernames must be at least 3 characters
// Usernames can only include latin letters, numbers, and one of -, _, or . but no special characters!
// Usernames cannot be longer than 15 characters
// Usernames can only include one -, _, or . You've got too many
// Usernames must start with a letter
var snapNameRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*[-_.]?[a-zA-Z0-9]+$`)

func (i *SnapchatChecker) ValidateName(name string) error {
	if len(name) < 3 {
		return fmt.Errorf("Name must be at least 3 characters")
	}
	if len(name) > 15 {
		return fmt.Errorf("Name must be less than 16 characters")
	}
	if !snapNameRe.MatchString(name) {
		return fmt.Errorf("Name must start with a letter and contain only letters, numbers, and at most one hyphen, underscore, or dot")
	}

	return nil
}

func (i *SnapchatChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://www.snapchat.com/add/"+name, nil)
	if status != 0 {
		return status
	}

	if strings.Contains(string(resp.Body), "<title data-react-helmet=\"true\">Snapchat</title>") {
		return StatusFree
	}

	return StatusUsed
}
