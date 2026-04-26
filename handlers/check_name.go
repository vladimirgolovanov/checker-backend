package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/vladimirgolovanov/checker-backend/namespaces"
)

type Request struct {
	Name       string             `json:"name"`
	Namespaces []NamespaceRequest `json:"namespaces"`
}

type NamespaceRequest struct {
	ID     int                    `json:"id"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type pendingItem struct {
	NamespaceID int    `json:"namespace_id"`
	Result      int    `json:"result"`
	Params      string `json:"params,omitempty"`
}

type resultItem struct {
	NamespaceID int    `json:"namespace_id"`
	Result      int    `json:"result"`
	Params      string `json:"params,omitempty"`
}

type validationErrorItem struct {
	NamespaceID int    `json:"namespace_id"`
	Errors      string `json:"errors"`
}

type checkerTask struct {
	checker      namespaces.Checker
	ns           NamespaceRequest
	preparedName string
	params       string
}

func writeSSEEvent(w http.ResponseWriter, event string, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		return
	}
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, b)
}

func CheckNameHandler(registry map[int]func(map[string]interface{}) []namespaces.Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}
		name := strings.ToLower(req.Name)

		// Pass 1: build pending list and task list sequentially
		var pending []pendingItem
		var tasks []checkerTask

		for _, ns := range req.Namespaces {
			factory, ok := registry[ns.ID]
			if !ok {
				continue
			}
			checkers := factory(ns.Params)
			for _, checker := range checkers {
				name = checker.PrepareName(name)
				item := pendingItem{
					NamespaceID: ns.ID,
					Result:      int(namespaces.StatusPending),
				}
				if dc, ok := checker.(*namespaces.DomainChecker); ok {
					item.Params = dc.Zone
				}
				pending = append(pending, item)
				tasks = append(tasks, checkerTask{
					checker:      checker,
					ns:           ns,
					preparedName: name,
					params:       item.Params,
				})
			}
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")

		writeSSEEvent(w, "pending", pending)
		flusher.Flush()

		// Pass 2: validate, send errors, launch goroutines
		resultCh := make(chan resultItem, len(tasks))
		var wg sync.WaitGroup

		for _, task := range tasks {
			if err := task.checker.ValidateName(task.preparedName); err != nil {
				writeSSEEvent(w, "validation_error", validationErrorItem{
					NamespaceID: task.ns.ID,
					Errors:      err.Error(),
				})
				flusher.Flush()
				continue
			}
			wg.Add(1)
			go func(t checkerTask) {
				defer wg.Done()
				status := t.checker.Check(t.preparedName, t.ns.Params)
				resultCh <- resultItem{
					NamespaceID: t.ns.ID,
					Result:      int(status),
					Params:      t.params,
				}
			}(task)
		}

		go func() {
			wg.Wait()
			close(resultCh)
		}()

		for res := range resultCh {
			writeSSEEvent(w, "result", res)
			flusher.Flush()
		}

		writeSSEEvent(w, "done", struct{}{})
		flusher.Flush()
	}
}
