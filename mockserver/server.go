package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

// store holds all resources in memory, keyed by name.
type store struct {
	mu             sync.RWMutex
	pools          map[string]json.RawMessage
	monitors       map[string]json.RawMessage
	virtualServers map[string]json.RawMessage
}

var db = &store{
	pools:          make(map[string]json.RawMessage),
	monitors:       make(map[string]json.RawMessage),
	virtualServers: make(map[string]json.RawMessage),
}

func main() {
	http.HandleFunc("/api/pools/", handleResource(db.pools, &db.mu))
	http.HandleFunc("/api/monitors/", handleResource(db.monitors, &db.mu))
	http.HandleFunc("/api/virtualservers/", handleResource(db.virtualServers, &db.mu))

	http.ListenAndServe(":8080", nil)
}

func handleResource(store map[string]json.RawMessage, mu *sync.RWMutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract resource name from URL path
		parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
		name := ""
		if len(parts) > 0 {
			name = parts[len(parts)-1]
		}

		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodGet:
			mu.RLock()
			data, ok := store[name]
			mu.RUnlock()
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Write(data)

		case http.MethodPost:
			var body json.RawMessage
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			mu.Lock()
			store[name] = body
			mu.Unlock()
			w.WriteHeader(http.StatusCreated)
			w.Write(body)

		case http.MethodPut:
			var body json.RawMessage
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			mu.Lock()
			store[name] = body
			mu.Unlock()
			w.Write(body)

		case http.MethodDelete:
			mu.Lock()
			delete(store, name)
			mu.Unlock()
			w.WriteHeader(http.StatusNoContent)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
