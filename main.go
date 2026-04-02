package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"

	"aboo.ru/checkers/namespaces"
)

type NamespaceRequest struct {
	ID     int                    `json:"id"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type Request struct {
	Name       string             `json:"name"`
	Namespaces []NamespaceRequest `json:"namespaces"`
}

type Namespaces struct {
	Namespace int                    `json:"namespace_id"`
	Result    namespaces.CheckStatus `json:"result"`
}

type Response struct {
	Results          []Namespaces             `json:"results"`
	ValidationErrors []map[string]interface{} `json:"validation_errors"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	sentryDsn := os.Getenv("SENTRY_DSN")
	if sentryDsn != "" {
		fmt.Printf("SENTRY_DSN: %s\n", sentryDsn)
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDsn,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
	}

	appHost := os.Getenv("APP_HOST")
	if appHost == "" {
		appHost = "0.0.0.0:8080"
	}

	http.HandleFunc("/api/check_name", checkNameHandler)
	log.Fatal(http.ListenAndServe(appHost, nil))
}

func checkNameHandler(w http.ResponseWriter, r *http.Request) {
	var checkerRegistry = map[int]func(params map[string]interface{}) []namespaces.Checker{
		0: single(&namespaces.InstagramChecker{}),
		1: func(params map[string]interface{}) []namespaces.Checker {
			zones := []string{"com"}
			if params != nil {
				if z, ok := params["zones"]; ok {
					if zoneSlice, ok := z.([]interface{}); ok {
						zones = make([]string, len(zoneSlice))
						for i, zone := range zoneSlice {
							zones[i] = zone.(string)
						}
					}
				}
			}
			checkers := make([]namespaces.Checker, len(zones))
			for i, zone := range zones {
				checkers[i] = &namespaces.DomainChecker{Zone: zone}
			}
			return checkers
		},
		5:  single(&namespaces.TiktokChecker{}),
		6:  single(&namespaces.SnapchatChecker{}),
		7:  single(&namespaces.NpmChecker{}),
		8:  single(&namespaces.GithubChecker{}),
		9:  single(&namespaces.TelegramChecker{}),
		10: single(&namespaces.TelegramBotChecker{}),
		11: single(&namespaces.EstyChecker{}),
		12: single(&namespaces.PinterestChecker{}),
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	name := strings.ToLower(req.Name)

	wg := sync.WaitGroup{}

	ch := make(chan Namespaces, len(req.Namespaces))
	var results []Namespaces
	var validationErrors []map[string]interface{}
	for _, ns := range req.Namespaces {
		factory, ok := checkerRegistry[ns.ID]
		if !ok {
			continue
		}

		checkers := factory(ns.Params)

		for _, checker := range checkers {
			if err := checker.ValidateName(name); err != nil {
				validationErrors = append(validationErrors, map[string]interface{}{
					"namespace": checker.GetId(),
					"errors":    err.Error(),
				})
				continue
			}

			wg.Add(1)
			go func(ns NamespaceRequest, checker namespaces.Checker) {
				checkerResult := checker.Check(name, ns.Params)
				ch <- Namespaces{
					Namespace: ns.ID,
					Result:    checkerResult,
				}
				wg.Done()
			}(ns, checker)
		}
	}

	wg.Wait()
	close(ch)

	for result := range ch {
		results = append(results, result)
	}

	fmt.Println(results)
	response := Response{
		Results:          results,
		ValidationErrors: validationErrors,
	}
	responseJSON, err := json.Marshal(response)
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

func single(c namespaces.Checker) func(map[string]interface{}) []namespaces.Checker {
	return func(params map[string]interface{}) []namespaces.Checker {
		return []namespaces.Checker{c}
	}
}
