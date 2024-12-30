package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	mu             sync.Mutex
	notifications  = make([]chan string, 0) // List of notification channels for each client
	submittedNames = make(chan string)      // Channel for submitted names
)

func eventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	tokens := []string{"foo", "bar", "baz", "qux", "quux", "corge", "grault", "garply", "waldo", "fred", "plugh", "xyzzy", "thud"}

	for _, token := range tokens {
		content := fmt.Sprintf("data: %s\n\n", string(token))
		w.Write([]byte(content))
		w.(http.Flusher).Flush()

		time.Sleep(time.Microsecond * 720)

	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// Notification stream handler
func notificationsHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create a notification channel for this client
	clientChan := make(chan string)
	mu.Lock()
	notifications = append(notifications, clientChan)
	mu.Unlock()

	// Remove the client's channel when they disconnect
	defer func() {
		mu.Lock()
		for i, ch := range notifications {
			if ch == clientChan {
				notifications = append(notifications[:i], notifications[i+1:]...)
				break
			}
		}
		mu.Unlock()
		close(clientChan)
	}()

	// Send notifications to the client
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	for notification := range clientChan {
		_, err := w.Write([]byte(fmt.Sprintf("data: %s\n\n", notification)))
		if err != nil {
			log.Println("Connection closed by client:", err)
			return
		}
		flusher.Flush()
	}
}

// Name submission handler
func submitNameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Send the submitted name to the notification channel
	go func() {
		submittedNames <- name
	}()
	fmt.Fprintf(w, "Name '%s' submitted successfully!", name)
}

// Broadcast submitted names to all clients
func broadcastNames() {
	for {
		name := <-submittedNames
		message := fmt.Sprintf("New name submitted: %s", name)

		mu.Lock()
		for _, clientChan := range notifications {
			clientChan <- message
		}
		mu.Unlock()
	}
}

func main() {
	// Start the broadcaster goroutine
	go broadcastNames()

	http.HandleFunc("/", handler)
	http.HandleFunc("/events", eventHandler)
	http.HandleFunc("/notifications", notificationsHandler)
	http.HandleFunc("/submit-name", submitNameHandler)

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
