package namespaces

import (
	"fmt"
	"regexp"
	"strings"
)

type YoutubeChecker struct{}

func (i *YoutubeChecker) GetId() int {
	return 13
}

func (i *YoutubeChecker) GetName() string {
	return "Youtube"
}

func (i *YoutubeChecker) GetSeparateSymbols() []string {
	return []string{
		"-",
		"_",
		".",
	}
}

func (i *YoutubeChecker) PrepareName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// a-z, 0-9, dot, hyphen, underscore; may not start or end with a special symbol.
var youtubeNameRe = regexp.MustCompile(`^[a-z0-9]([a-z0-9._-]*[a-z0-9])?$`)

// two consecutive special symbols (e.g. .., --, __, .-)
var youtubeConsecutiveSpecialRe = regexp.MustCompile(`[._-]{2}`)

func (i *YoutubeChecker) ValidateName(name string) error {
	if len(name) < 3 || len(name) > 30 {
		return fmt.Errorf("Name must be between 3 and 30 characters")
	}

	if !youtubeNameRe.MatchString(name) {
		return fmt.Errorf("Name must contain only letters, numbers, dots, hyphens, or underscores, and cannot begin or end with a dot, hyphen, or underscore")
	}

	if youtubeConsecutiveSpecialRe.MatchString(name) {
		return fmt.Errorf("Name must not contain two consecutive special characters")
	}

	return nil
}

func (i *YoutubeChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://www.youtube.com/@"+name, nil)
	if status != 0 {
		return status
	}

	if resp.StatusCode != 404 {
		return StatusUsed
	}

	return StatusFree
}
