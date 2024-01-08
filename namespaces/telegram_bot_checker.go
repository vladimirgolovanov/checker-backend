package namespaces

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type TelegramBotChecker struct {
}

func (i *TelegramBotChecker) GetId() int {
	return 10
}

func (i *TelegramBotChecker) Check(name string) CheckStatus {
	url := "https://t.me/" + name + "bot"
	response, err := http.Get(url)
	if err != nil {
		return StatusFailed
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return StatusFailed
	}

	if strings.Contains(string(body), "<meta property=\"twitter:title\" content=\"Telegram: Contact @") {
		return StatusFree
	}

	return StatusUsed
}
