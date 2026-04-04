package namespaces

import (
	"fmt"
	"strings"
)

type InstagramChecker struct{}

func (i *InstagramChecker) GetId() int {
	return 0
}

func (i *InstagramChecker) GetName() string {
	return "Instagram"
}

func (i *InstagramChecker) PrepareName(name string) string {
	return name
}

func (i *InstagramChecker) ValidateName(name string) error {
	return nil
}

func (i *InstagramChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://www.instagram.com/"+name+"/", nil)
	if status != 0 {
		return status
	}

	body := string(resp.Body)

	if strings.Contains(body, "{\"username\":\""+name+"\"}") {
		fmt.Println("used")
		return StatusUsed
	}

	if !strings.Contains(body, "\"url\":\"\\/"+name+"\\/\"") {
		fmt.Println("failed on finding url string")
		return StatusFailed
	}

	fmt.Println("free")
	return StatusFree
}
