package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type NetDomainChecker struct{}

func (i *NetDomainChecker) GetId() int {
	return 3
}

func (i *NetDomainChecker) Check(name string) bool {
	domainName := name + ".net"
	cmd := exec.Command("whois", domainName)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Ошибка хз чего:", err)
	}

	if strings.Contains(string(output), "No match for domain") {
		return true
	}

	return false
}
