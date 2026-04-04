package namespaces

type EtsyChecker struct{}

func (i *EtsyChecker) GetId() int {
	return 11
}

func (i *EtsyChecker) GetName() string {
	return "Esty"
}

func (i *EtsyChecker) PrepareName(name string) string {
	return name
}

// если содержит a-z, 0-9, _ и длиной от 4 до 20 символов
// если нет, то возвращаем ошибку
func (i *EtsyChecker) ValidateName(name string) error {
	return nil
}

func (i *EtsyChecker) Check(name string, params map[string]interface{}) CheckStatus {
	resp, status := Get("https://www.etsy.com/shop/"+name, nil)
	if status != 0 {
		return status
	}

	if resp.StatusCode == 404 {
		return StatusFree
	}

	return StatusUsed
}
