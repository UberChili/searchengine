package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ApiResponse represents the top-level response from Amadeus API
type ApiResponse struct {
  Data []flight `json:"data"`
}

// Price represents the price information of a flight
type Price struct {
  Total string `json:"total"`
}

// flight represents data about a flight
type flight struct {
  Origin string  `json:"origin"`
  Destination string  `json:"destination"`
  // Airline string  `json:"airline"`
  // FlightNumber string  `json:"flight_number"`
  Price  Price `json:"price"`
}

const (
  baseURL = "https://test.api.amadeus.com/v1/shopping/flight-destinations"
  authToken = "Bearer hLfK0atgerjfXZoFAzqMmpb6DxFy"
)

func main() {
  // Amadeus secrets
  var API_KEY = os.Getenv("AMAD_API_KEY")
  var SECRET = os.Getenv("AMAD_SECRET")

  router := gin.Default()
  router.GET("/flights/:origin", getFlights)

  router.Run("localhost:8080")
}

// getFlights responds with the list of all flights from a destination as JSON
func getFlights(c *gin.Context) {
  origin := c.Param("origin")

  // Create the request
  req, err := http.NewRequest("GET", fmt.Sprintf(
    "%s?origin=%s&maxPrice=200", baseURL, origin), nil)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  // Add headers
  req.Header.Add("Authorization", authToken)

  // Make the request
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  defer resp.Body.Close()

  // Read response Body
  body, err := io.ReadAll(resp.Body)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  // Check if the response was successful
  if resp.StatusCode != http.StatusOK {
    c.JSON(resp.StatusCode, gin.H{"error": string(body)})
    return
  }

  // Parse the JSON response
  var response ApiResponse
  if err := json.Unmarshal(body, &response); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response: " + err.Error()})
    return
  }

  c.IndentedJSON(http.StatusFound, response.Data)
}
