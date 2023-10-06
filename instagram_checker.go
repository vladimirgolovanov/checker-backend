package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type InstagramChecker struct{}

func (i *InstagramChecker) GetId() int {
	return 0
}

func (i *InstagramChecker) Check(name string) bool {
	url := "http://localhost:3000/api/v1/check-username"
	data := map[string]interface{}{
		"username": name,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Ошибка при кодировании JSON данных:", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		// временно забьем на ошибку
		fmt.Println("Ошибка при выполнении запроса:", err)
	}
	req.Header.Set("Authorization", "your_auth_token_here")
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке HTTP запроса:", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("Ошибка при чтении JSON ответа:", err)
	}

	return response["success"] == true
}
