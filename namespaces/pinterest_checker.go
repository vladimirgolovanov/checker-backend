package namespaces

import "strings"

type PinterestChecker struct{}

func (i *PinterestChecker) GetId() int {
	return 12
}

func (i *PinterestChecker) GetName() string {
	return "Pinterest"
}

func (i *PinterestChecker) PrepareName(name string) string {
	return name
}

func (i *PinterestChecker) ValidateName(name string) error {
	return nil
}

func (i *PinterestChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://www.pinterest.com/"+name+"/", nil)
	if status != 0 {
		return status
	}

	if strings.Contains(string(resp.Body), "Profile | Pinterest</title>") {
		return StatusUsed
	}

	return StatusFree
}
