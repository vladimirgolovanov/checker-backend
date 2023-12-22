package namespaces

import (
	"os/exec"
	"strings"
)

type RuDomainChecker struct{}

func (i *RuDomainChecker) GetId() int {
	return 2
}

func (i *RuDomainChecker) Check(name string) CheckStatus {
	domainName := name + ".ru"
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
