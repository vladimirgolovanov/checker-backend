package main

import (
	"aboo.ru/checkers/namespaces"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type Request struct {
	Name       string `json:"name"`
	Namespaces []int  `json:"namespaces"`
}

type Namespaces struct {
	Namespace int                    `json:"namespace_id"`
	Result    namespaces.CheckStatus `json:"result"`
}

func main() {
	http.HandleFunc("/", checkNames)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func checkNames(w http.ResponseWriter, r *http.Request) {
	services := []namespaces.Checker{
		&namespaces.InstagramChecker{}, // 0
		&namespaces.ComDomainChecker{}, // 1
		&namespaces.RuDomainChecker{},  // 2
		&namespaces.NetDomainChecker{}, // 3
		&namespaces.IoDomainChecker{},  // 4
		&namespaces.TiktokChecker{},    // 5
		&namespaces.SnapchatChecker{},  // 6
		// npm 7
		&namespaces.GithubChecker{},      // 8
		&namespaces.TelegramChecker{},    // 9
		&namespaces.TelegramBotChecker{}, // 10
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData Request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Ошибка при закрытии Body.Close:", err)
		}
	}(r.Body)

	name := requestData.Name
	var filteredServices []namespaces.Checker

	for _, service := range services {
		for _, namespace := range requestData.Namespaces {
			if service.GetId() == namespace {
				filteredServices = append(filteredServices, service)
				break
			}
		}
	}

	var wg = sync.WaitGroup{}

	fmt.Println(filteredServices)

	ch := make(chan Namespaces, len(filteredServices))
	var results []Namespaces
	for _, service := range filteredServices {
		wg.Add(1)
		service := service // Loop variables captured by 'func' literals in 'go' statements might have unexpected values
		go func() {
			serviceResult := service.Check(name)
			newNamespace := Namespaces{
				Namespace: service.GetId(),
				Result:    serviceResult,
			}
			ch <- newNamespace
			wg.Done()
		}()
	}

	go func() {
		for {
			result := <-ch
			results = append(results, result)
		}
	}()

	wg.Wait()

	fmt.Println(results)
	responseJSON, err := json.Marshal(results)
	if err != nil {
		fmt.Println("Ошибка при маршалинге JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		fmt.Println("Ошибка отправки ответа:", err)
	}
}
