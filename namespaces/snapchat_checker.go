package namespaces

import "strings"

type SnapchatChecker struct{}

func (i *SnapchatChecker) GetId() int {
	return 6
}

func (i *SnapchatChecker) GetName() string {
	return "Snapchat"
}

func (i *SnapchatChecker) PrepareName(name string) string {
	return name
}

func (i *SnapchatChecker) ValidateName(name string) error {
	return nil
}

func (i *SnapchatChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://www.snapchat.com/add/"+name, nil)
	if status != 0 {
		return status
	}

	if strings.Contains(string(resp.Body), "<title data-react-helmet=\"true\">Snapchat</title>") {
		return StatusFree
	}

	return StatusUsed
}
