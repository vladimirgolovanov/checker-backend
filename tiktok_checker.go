package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type TiktokChecker struct{}

func (i *TiktokChecker) GetId() int {
	return 5
}

func (i *TiktokChecker) Check(name string) bool {
	url := "https://www.tiktok.com/@" + name
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

	return strings.Contains(string(body), "uniqueId")
}
