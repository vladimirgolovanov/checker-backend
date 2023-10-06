package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type IoDomainChecker struct{}

func (i *IoDomainChecker) GetId() int {
	return 4
}

func (i *IoDomainChecker) Check(name string) bool {
	domainName := name + ".io"
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
