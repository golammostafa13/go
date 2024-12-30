package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Notification structure
type Notification struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Name storage (mock database)
var submittedNames = []string{}
var notifications = []Notification{}
var notificationsLock sync.Mutex

// WebSocket clients and broadcaster
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Notification)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket handler
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket Upgrade Failed", http.StatusInternalServerError)
		return
	}
	defer ws.Close()
	clients[ws] = true

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			break
		}
	}
}

// Handle broadcasts to WebSocket clients
func handleBroadcasts() {
	for {
		notification := <-broadcast
		for client := range clients {
			err := client.WriteJSON(notification)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// Submit name and create notification
func submitNameHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Add the name to the submitted names (mock database)
	submittedNames = append(submittedNames, name)

	// Create and broadcast a notification
	notification := Notification{
		ID:        len(notifications) + 1,
		Content:   fmt.Sprintf("New name submitted: %s", name),
		Timestamp: time.Now(),
	}

	notificationsLock.Lock()
	notifications = append(notifications, notification)
	notificationsLock.Unlock()

	broadcast <- notification

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Name submitted successfully!"})
}

// Fetch notifications (historical data)
func getNotifications(w http.ResponseWriter, r *http.Request) {
	notificationsLock.Lock()
	defer notificationsLock.Unlock()
	json.NewEncoder(w).Encode(notifications)
}

func handle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	http.HandleFunc("/", handle)
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/submit-name", submitNameHandler)
	http.HandleFunc("/notifications", getNotifications)

	go handleBroadcasts()

	fmt.Println("Server is running on http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
