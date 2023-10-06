package main

import (
	"fmt"
	"net/http"
)

type GithubChecker struct{}

func (i *GithubChecker) GetId() int {
	return 8
}

func (i *GithubChecker) Check(name string) bool {
	url := "https://github.com/" + name
	response, err := http.Get(url)
	if err != nil {
		// временно забьем на ошибку
		fmt.Println("Ошибка при выполнении запроса:", err)
	}
	defer response.Body.Close()

	return response.StatusCode != 404
}
