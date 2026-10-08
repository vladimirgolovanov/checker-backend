package namespaces

import (
	"fmt"
	"regexp"
	"strings"
)

type GithubChecker struct{}

func (i *GithubChecker) GetId() int {
	return 8
}

func (i *GithubChecker) GetName() string {
	return "Github"
}

func (i *GithubChecker) GetSeparateSymbols() []string {
	return []string{}
}

func (i *GithubChecker) PrepareName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// Username may only contain alphanumeric characters or single hyphens, and cannot begin or end with a hyphen.
// maximum is 39 characters
var githubNameRe = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

func (i *GithubChecker) ValidateName(name string) error {
	if len(name) > 39 {
		return fmt.Errorf("Name must be less than 40 characters")
	}

	if !githubNameRe.MatchString(name) {
		return fmt.Errorf("Name must contain only letters, numbers, or hyphens, and cannot begin or end with a hyphen")
	}

	if strings.Contains(name, "--") {
		return fmt.Errorf("Name must not contain consecutive hyphens")
	}

	return nil
}

func (i *GithubChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://github.com/"+name, nil)
	if status != 0 {
		return status
	}

	if resp.StatusCode != 404 {
		return StatusUsed
	}

	return StatusFree
}
