package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
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
		&SnapchatChecker{},  // 6
		&ComDomainChecker{}, // 1
		&RuDomainChecker{},  // 2
		&NetDomainChecker{}, // 3
		&IoDomainChecker{},
		&GithubChecker{},
	}

	//log.Println("r.URL.Path")
	//log.Println(r.URL.Path)

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

	//log.Println(requestData.Name)
	//log.Println(requestData.Namespaces)
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

	//log.Println(filteredServices)

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

	timeout := time.After(2 * time.Second)

	go func() {
		for {
			select {
			case result := <-ch:
				results = append(results, result)
				log.Println("new results")
			case <-timeout:
				close(ch)
			}
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

//func responseResults(w http.ResponseWriter, results []Namespaces) {
//	log.Println(len(results))
//
//	responseJSON, err := json.Marshal(results)
//	if err != nil {
//		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	w.Write(responseJSON)
//
//	log.Println(responseJSON)
//}
