package namespaces

type GithubChecker struct{}

func (i *GithubChecker) GetId() int {
	return 8
}

func (i *GithubChecker) GetName() string {
	return "Github"
}

func (i *GithubChecker) PrepareName(name string) string {
	return name
}

func (i *GithubChecker) ValidateName(name string) error {
	return nil
}

func (i *GithubChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://github.com/"+name, nil)
	if status != 0 {
		return status
	}

	if resp.StatusCode != 404 {
		return StatusUsed
	}

	return StatusFree
}
