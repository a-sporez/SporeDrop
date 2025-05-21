// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// incoming message structure
type Message struct {
	Message string `json:"message"`
}

// outgoing reply structure
type Response struct {
	Reply string `json:"reply"`
}

func chat_handler(w http.ResponseWriter, r *http.Request) {
	var msg Message

	// decode incoming json
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userMsg := strings.ToLower(msg.Message)
	reply := "I'm not sure what you mean by that."

	// keyword based logic
	switch {
	case strings.Contains(userMsg, "hello"):
		reply = "Hello! How can I help you?"
	case strings.Contains(userMsg, "bye"):
		reply = "Goodbye!"
	case strings.Contains(userMsg, "help"):
		reply = "Yes! I am here to assist you!"
	}

	// send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Reply: reply})
}

func main() {
	http.HandleFunc("/chat", chat_handler)
	log.Println("chatbot running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
