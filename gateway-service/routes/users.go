package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"shopsweb.com/gateway-service/models"
)

var authServiceURL = os.Getenv("AUTH_SERVICE_URL")
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func signup(context *gin.Context) {
	var user models.User

	// decided to separate this from the login so i can change parameters from one another
	if err := context.ShouldBindJSON(&user); err != nil {
		log.Printf("Error binding JSON: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data. Ensure email, and password are provided."})
		return
	}

	handleAuthRequest(context, "signup", &user)
}

func login(context *gin.Context) {
	var user models.User

	// decided to separate this from the signup so i can change parameters from one another
	if err := context.ShouldBindJSON(&user); err != nil {
		log.Printf("Error binding JSON: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data. Ensure email, and password are provided."})
		return
	}

	handleAuthRequest(context, "login", &user)
}

func handleAuthRequest(context *gin.Context, endpoint string, user *models.User) {
	// Serialize User Data for the Auth Service
	jsonData, err := json.Marshal(&user)
	if err != nil {
		log.Printf("Error marshalling user data: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not process user data."})
		return
	}

	// Create Request to Auth Service
	req, err := http.NewRequest("POST", authServiceURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request to auth service: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create authentication request."})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send Request to Auth Service
	log.Printf("Sending %s request for user %s to %s", endpoint, user.Email, authServiceURL)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error sending request to auth service at %s: %v", authServiceURL, err)
		context.JSON(http.StatusServiceUnavailable, gin.H{"message": fmt.Sprintf("Could not reach authentication service: %v", err)})
		return
	}
	defer resp.Body.Close()

	log.Printf("Received response from auth service: Status %d", resp.StatusCode)

	// Handle Response from Auth Service
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var authResponse map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&authResponse); err == nil {
			log.Printf("Auth service response body: %v", authResponse)
			context.JSON(resp.StatusCode, authResponse) // Forward the response
		} else {
			// If parsing fails but status is success, send a generic success
			context.JSON(resp.StatusCode, gin.H{"message": "Operation completed successfully."})
		}
		return
	}

	// Handle error response
	var errorResponse gin.H
	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil && errorResponse["message"] != nil {
		log.Printf("Error from auth service (%d): %v", resp.StatusCode, errorResponse["message"])
		context.JSON(resp.StatusCode, errorResponse)
	} else {
		log.Printf("Error from auth service (%d), but could not parse error body.", resp.StatusCode)
		context.JSON(resp.StatusCode, gin.H{"message": fmt.Sprintf("Authentication service returned status %d.", resp.StatusCode)})
	}
}
