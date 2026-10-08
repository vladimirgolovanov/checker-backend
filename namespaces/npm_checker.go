package namespaces

import (
	"fmt"
	"regexp"
	"strings"
)

type NpmChecker struct{}

func (i *NpmChecker) GetId() int {
	return 7
}

func (i *NpmChecker) GetName() string {
	return "Npm"
}

func (i *NpmChecker) GetSeparateSymbols() []string {
	return []string{
		"-",
		"_",
		".",
	}
}

func (i *NpmChecker) PrepareName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// Name length must be less than or equal to 214 characters long
// Name may not start with "."
// Name must be lowercase
var npmNameRe = regexp.MustCompile(`^[a-z0-9._~()!*-]+$`) // todo: real check of ~()!* symbols

func (i *NpmChecker) ValidateName(name string) error {
	if len(name) > 214 {
		return fmt.Errorf("Name must be less than 215 characters")
	}
	if strings.HasPrefix(name, ".") {
		return fmt.Errorf("Name must not start with a dot")
	}
	if !npmNameRe.MatchString(name) {
		return fmt.Errorf("Name must contain only lowercase letters, numbers, hyphens, underscores, or dots")
	}

	return nil
}

func (i *NpmChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://www.npmjs.com/package/"+name, nil)
	if status != 0 {
		return status
	}

	if resp.StatusCode == 404 {
		return StatusFree
	}

	return StatusUsed
}
