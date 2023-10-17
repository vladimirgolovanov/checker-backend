package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type RuDomainChecker struct{}

func (i *RuDomainChecker) GetId() int {
	return 2
}

func (i *RuDomainChecker) Check(name string) bool {
	domainName := name + ".ru"
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
