package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// PaymentEvent is the struct for boxo payment callback
type PaymentEvent struct {
	TransactionToken string          `json:"transactionToken"`
	MiniappOrderId   int64           `json:"miniappOrderId"`
	HostappOrderId   int64           `json:"hostappOrderId"`
	Amount           float64         `json:"amount"`
	Currency         string          `json:"currency"`
	Status           string          `json:"status"`
	ExtraParams      json.RawMessage `json:"extraParams"`
}

func paymentWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Validating Basic Auth
	auth := r.Header.Get("Authorization")
	expectedUser := os.Getenv("CLIENT_ID")
	expectedPass := os.Getenv("CLIENT_SECRET")

	if auth == "" || expectedUser == "" || expectedPass == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	expected := "Basic " + base64.StdEncoding.EncodeToString([]byte(expectedUser+":"+expectedPass))
	if auth != expected {
		log.Printf("Invalid auth header: %s", auth)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Body read error: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var evt PaymentEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("JSON parse error: %v – %s", err, string(body))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Log the event
	log.Printf("Payment received: tx=%s, miniappOrder=%d, hostappOrder=%d, amount=%.2f %s, status=%s",
		evt.TransactionToken, evt.MiniappOrderId, evt.HostappOrderId, evt.Amount, evt.Currency, evt.Status)
	log.Printf("ExtraParams: %s", string(evt.ExtraParams))

	// Add business logic here

	// Respond OK to stop retries
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"received"}`))
}

// NOTE: wrote this with a separate main so we have the ws logic seperated
// because we haven't implemented the complete integration
func main() {
	r := mux.NewRouter()
	// NOTE: Using sample webhook endpoint
	r.HandleFunc("/webhook/boxo/payment", paymentWebhookHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on :%s/webhook/boxo/payment", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
