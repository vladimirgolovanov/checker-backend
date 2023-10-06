package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type ComDomainChecker struct{}

func (i *ComDomainChecker) GetId() int {
	return 1
}

func (i *ComDomainChecker) Check(name string) bool {
	return true
	domainName := name + ".com"
	cmd := exec.Command("whois", domainName)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Ошибка хз чего:", err)
	}

	if strings.Contains(string(output), "No match for domain") {
		return true
	}

	return false

	//fmt.Println("Результат whois для домена", domain_name, ":\n", string(output))

	//return true
}
