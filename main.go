package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type TokenRequest struct {
	AuthCode string `json:"authCode"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	ErrorCode   string `json:"error_code,omitempty"`
}

func getAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	serverURL := os.Getenv("AUTH_SERVER_URL")

	if clientID == "" || clientSecret == "" || serverURL == "" {
		http.Error(w, "Missing environment variables", http.StatusInternalServerError)
		return
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientID, clientSecret)))
	payload, _ := json.Marshal(map[string]string{"authCode": req.AuthCode})

	httpReq, err := http.NewRequest("POST", serverURL+"/api/get_access_token/", bytes.NewBuffer(payload))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, "Failed to reach auth server", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		http.Error(w, "Invalid response from auth server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenResp)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/get-token", getAccessTokenHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
