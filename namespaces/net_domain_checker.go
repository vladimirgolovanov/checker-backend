package namespaces

import (
	"os/exec"
	"strings"
)

type NetDomainChecker struct{}

func (i *NetDomainChecker) GetId() int {
	return 3
}

func (i *NetDomainChecker) Check(name string) CheckStatus {
	domainName := name + ".net"
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
