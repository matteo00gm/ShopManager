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

	// Parse Incoming Request
	err := context.ShouldBindJSON(&user)
	if err != nil {
		log.Printf("Error binding JSON: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data. Ensure username, email, and password are provided."})
		return
	}

	// Serialize User Data for the Auth Service
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling user data: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not process user data."})
		return
	}

	// Create Request to Auth Service
	req, err := http.NewRequest("POST", authServiceURL+"signup", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request to auth service: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create authentication request."})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send Request to Auth Service
	log.Printf("Sending signup request for user %s to %s", user.Email, authServiceURL)
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
			context.JSON(resp.StatusCode, gin.H{"message": "User created successfully."})
		}
		return
	}

	// error
	var errorResponse gin.H
	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil && errorResponse["message"] != nil {
		log.Printf("Error from auth service (%d): %v", resp.StatusCode, errorResponse["message"])
		// Forward the error structure if possible, using the status code from the auth service
		context.JSON(resp.StatusCode, errorResponse)
	} else {
		// Fallback if body parsing fails or doesn't contain a 'message' field
		log.Printf("Error from auth service (%d), but could not parse error body.", resp.StatusCode)
		context.JSON(resp.StatusCode, gin.H{"message": fmt.Sprintf("Authentication service returned status %d.", resp.StatusCode)})
	}
}
