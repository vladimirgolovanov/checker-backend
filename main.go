package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Request struct {
	Name       string `json:"name"`
	Namespaces []int  `json:"namespaces"`
}

type Namespaces struct {
	Namespace int  `json:"namespace_id"`
	Result    bool `json:"result"`
}

func main() {
	http.HandleFunc("/", checkNames)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func checkNames(w http.ResponseWriter, r *http.Request) {
	services := []Checker{
		&InstagramChecker{}, // 0
		&ComDomainChecker{}, // 1
		&RuDomainChecker{},  // 2
		&NetDomainChecker{}, // 3
		&IoDomainChecker{},  // 4
		&SnapchatChecker{},  // 6
		&GithubChecker{},    // 8
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
	defer r.Body.Close()

	name := requestData.Name
	var filteredServices []Checker

	for _, service := range services {
		for _, namespace := range requestData.Namespaces {
			if service.GetId() == namespace {
				filteredServices = append(filteredServices, service)
				break
			}
		}
	}

	var wg = sync.WaitGroup{}

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
	w.Write(responseJSON)
}
