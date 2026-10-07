package namespaces

import "strings"

// Usernames can't contain numbers only
// Only lowercase letters, numbers, underscores, or periods are allowed
// Usernames can't end with a period
// min 2, max 24

type TiktokChecker struct{}

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
	return nil
}

func (i *TiktokChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://www.tiktok.com/@"+name, nil)
	if status != 0 {
		return status
	}

	if strings.Contains(string(resp.Body), "uniqueId") {
		return StatusUsed
	}

	return StatusFree
}
