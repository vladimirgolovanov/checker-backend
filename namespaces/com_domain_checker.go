package namespaces

import (
	"os/exec"
	"strings"
)

type ComDomainChecker struct{}

func (i *ComDomainChecker) GetId() int {
	return 1
}

func (i *ComDomainChecker) Check(name string) CheckStatus {
	domainName := name + ".com"
	cmd := exec.Command("whois", domainName)

	output, err := cmd.Output()
	if err != nil {
		return StatusFailed
	}

	if strings.Contains(string(output), "No match for domain") {
		return StatusFree
	}

	return StatusUsed
}
