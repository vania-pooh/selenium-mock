package main

import (
	"net/http"
	"sync"
	"encoding/json"
	"strings"
	"flag"
	"log"
	"fmt"

	"github.com/pborman/uuid"
)

var (
	port int
)

const (
	screenshot = "iVBORw0KGgoAAAANSUhEUgAAASwAAADhCAMAAABFoniZAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyBpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuMC1jMDYwIDYxLjEzNDc3NywgMjAxMC8wMi8xMi0xNzozMjowMCAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENTNSBXaW5kb3dzIiB4bXBNTTpJbnN0YW5jZUlEPSJ4bXAuaWlkOjgwOTkzNjZBMzI5RTExRTA5RjU1Q0RGODVCQjc2ODMxIiB4bXBNTTpEb2N1bWVudElEPSJ4bXAuZGlkOjgwOTkzNjZCMzI5RTExRTA5RjU1Q0RGODVCQjc2ODMxIj4gPHhtcE1NOkRlcml2ZWRGcm9tIHN0UmVmOmluc3RhbmNlSUQ9InhtcC5paWQ6ODA5OTM2NjgzMjlFMTFFMDlGNTVDREY4NUJCNzY4MzEiIHN0UmVmOmRvY3VtZW50SUQ9InhtcC5kaWQ6ODA5OTM2NjkzMjlFMTFFMDlGNTVDREY4NUJCNzY4MzEiLz4gPC9yZGY6RGVzY3JpcHRpb24+IDwvcmRmOlJERj4gPC94OnhtcG1ldGE+IDw/eHBhY2tldCBlbmQ9InIiPz5JKQ/cAAAAMFBMVEWWlpZ3d3dZWVk7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"
)

func handler() http.Handler {
	var lock sync.RWMutex
	sessions := make(map[string]struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		u := uuid.New()
		lock.Lock()
		sessions[u] = struct{}{}
		lock.Unlock()
		json.NewEncoder(w).Encode(struct {
			S string `json:"sessionId"`
		}{u})
	})
	mux.HandleFunc("/session/", func(w http.ResponseWriter, r *http.Request) {
		u := strings.Split(r.URL.Path, "/")[2]
		lock.RLock()
		_, ok := sessions[u]
		lock.RUnlock()
		if !ok {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(struct {
			St string `json:"state"`
			Id string `json:"sessionId"`
			Value string `json:"value"`
		}{"success", u, screenshot})
		if r.Method != http.MethodDelete {
			return
		}
		lock.Lock()
		delete(sessions, u)
		lock.Unlock()
	})
	return mux
}

func init() {
	flag.IntVar(&port, "port", 4444, "port to bind to")
	flag.Parse()
}

func main() {
	listen := fmt.Sprintf(":%d", port)
	log.Println("listening on", listen)
	log.Print(http.ListenAndServe(listen, handler()))
}