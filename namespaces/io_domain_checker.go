package namespaces

import (
	"os/exec"
	"strings"
)

type IoDomainChecker struct{}

func (i *IoDomainChecker) GetId() int {
	return 4
}

func (i *IoDomainChecker) Check(name string) CheckStatus {
	domainName := name + ".io"
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
