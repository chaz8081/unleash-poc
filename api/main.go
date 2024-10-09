package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Unleash/unleash-client-go/v4"
)

var unleashClient *unleash.Client

func init() {
	// Get Unleash configuration from environment variables
	unleashURL := os.Getenv("UNLEASH_URL")
	appName := os.Getenv("UNLEASH_APP_NAME")
	instanceID := os.Getenv("UNLEASH_INSTANCE_ID")
	clientApiToken := os.Getenv("INIT_CLIENT_API_TOKENS")

	// Initialize Unleash client
	var err error
	unleashClient, err = unleash.NewClient(
		unleash.WithListener(&unleash.DebugListener{}),
		unleash.WithUrl(unleashURL),
		unleash.WithAppName(appName),
		unleash.WithInstanceId(instanceID), // Optional
		unleash.WithCustomHeaders(http.Header{"Authorization": {clientApiToken}}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Note this will block until the default client is ready
	unleashClient.WaitForReady()
}

func main() {
	defer unleashClient.Close()

	log.Printf("features: %+v", unleashClient.ListFeatures())

	router := http.NewServeMux()

	// API endpoint handler
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the feature flag is enabled
		if unleashClient.IsEnabled("new-feature") {
			fmt.Fprint(w, "new-feature is enabled!")
		} else {
			fmt.Fprint(w, "new-feature is disabled.")
		}
	})

	// Start the API server
	log.Println("API server listening on :8080")
	http.ListenAndServe(":8080", router)
}
