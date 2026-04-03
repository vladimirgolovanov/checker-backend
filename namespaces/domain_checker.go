package namespaces

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type DomainChecker struct {
	Zone string
}

func (i *DomainChecker) GetId() int {
	return 1
}

func (i *DomainChecker) GetName() string {
	return "Domain name"
}

func (i *DomainChecker) PrepareName(name string) string {
	return name
}

func (i *DomainChecker) ValidateName(name string) error {
	if len(name) < 2 {
		return errors.New("Name must be at least 2 characters long")
	}

	if len(name) > 63 {
		return errors.New("Name must not exceed 63 characters")
	}

	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return errors.New("Name may consist only of a-z, 0-9, and hyphens")
		}
	}

	if name[0] == '-' || name[len(name)-1] == '-' {
		return errors.New("Name must not start or end with a hyphen")
	}

	return nil
}

func (i *DomainChecker) Check(name string, params map[string]interface{}) CheckStatus {
	domainName := name + "." + i.Zone

	cmd := exec.Command("whois", domainName)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("failed on whois: ", err)
		return StatusFailed
	}

	if strings.Contains(string(output), "No match for domain") {
		return StatusFree
	}

	return StatusUsed
}
