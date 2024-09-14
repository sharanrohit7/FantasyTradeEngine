package API

import (
	"TradeEngine/common"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func MakeDebitRequest(ctx context.Context, transaction common.DebitTransaction) error {
	baseURL := os.Getenv("WALLET_URL")
	url := fmt.Sprintf("%s/debit-money", baseURL)

	// Create the request body
	requestBody, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("error marshaling request body: %w", err)
	}

	// Create the HTTP POST request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}
	log.Println("Debit request successful")

	return nil
}

type TokenResponse struct {
	UserId string `json:"id"`
	RoleId string `json:"role_id"`
}

func ExtractUserIdFromToken(token string) (string, error) {
	// Get the token microservice base URL from environment variable
	baseUrl := os.Getenv("TOKEN_URL")
	if baseUrl == "" {
		return "", errors.New("TOKEN_URL environment variable is not set")
	}

	// Construct the full URL
	url := fmt.Sprintf("%s/extractData", baseUrl)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add token to request header
	req.Header.Add("Authorization", token)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	// fmt.Println("TOKEN REPONSE BODY IS ..............................", resp)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check for valid response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to extract userId: status code %d", resp.StatusCode)
	}

	// Parse the response
	var tokenResp TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.UserId == "" {
		return "", errors.New("userId not found in token response")
	}

	return tokenResp.UserId, nil
}
