package namespaces

type NpmChecker struct{}

func (i *NpmChecker) GetId() int {
	return 7
}

func (i *NpmChecker) GetName() string {
	return "Npm"
}

func (i *NpmChecker) PrepareName(name string) string {
	return name
}

func (i *NpmChecker) ValidateName(name string) error {
	return nil
}

func (i *NpmChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://www.npmjs.com/package/"+name, nil)
	if status != 0 {
		return status
	}

	if resp.StatusCode == 404 {
		return StatusFree
	}

	return StatusUsed
}
