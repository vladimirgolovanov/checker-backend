package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"aboo.ru/checkers/namespaces"
)

type Request struct {
	Name       string             `json:"name"`
	Namespaces []NamespaceRequest `json:"namespaces"`
}

type NamespaceRequest struct {
	ID     int                    `json:"id"`
	Params map[string]interface{} `json:"params,omitempty"`
}
type Namespaces struct {
	Namespace int                    `json:"namespace_id"`
	Result    namespaces.CheckStatus `json:"result"`
	Params    string                 `json:"params,omitempty"`
}

type Response struct {
	Results          []Namespaces             `json:"results"`
	ValidationErrors []map[string]interface{} `json:"validation_errors"`
}

func CheckNameHandler(registry map[int]func(map[string]interface{}) []namespaces.Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var mu sync.Mutex
		var results []Namespaces
		var validationErrors []map[string]interface{}
		for _, ns := range req.Namespaces {
			factory, ok := registry[ns.ID]
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
					defer wg.Done()
					checkerResult := checker.Check(name, ns.Params)
					entry := Namespaces{Namespace: ns.ID, Result: checkerResult}
					if dc, ok := checker.(*namespaces.DomainChecker); ok {
						entry.Params = dc.Zone
					}
					mu.Lock()
					results = append(results, entry)
					mu.Unlock()
				}(ns, checker)
			}
		}

		wg.Wait()

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
}
