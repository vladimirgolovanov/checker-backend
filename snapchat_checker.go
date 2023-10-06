package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type SnapchatChecker struct{}

func (i *SnapchatChecker) GetId() int {
	return 6
}

func (i *SnapchatChecker) Check(name string) bool {
	url := "https://www.snapchat.com/add/" + name
	response, err := http.Get(url)
	if err != nil {
		// временно забьем на ошибку
		fmt.Println("Ошибка при выполнении запроса:", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// временно забьем на ошибку
		fmt.Println("Ошибка при выполнении запроса:", err)
	}

	return strings.Contains(string(body), "<title data-react-helmet=\"true\">Snapchat</title>")
}
