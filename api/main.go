package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Unleash/unleash-client-go/v4"
	uContext "github.com/Unleash/unleash-client-go/v4/context"
	"github.com/google/uuid"
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

	router := http.NewServeMux()

	// API endpoint handler
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the feature flag is enabled
		if unleashClient.IsEnabled("new-feature") {
			fmt.Fprintln(w, "new-feature is enabled!")
		} else {
			fmt.Fprintln(w, "new-feature is disabled.")
		}
	})

	router.HandleFunc("/buyer/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")

		// Parse UUID
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid buyer ID", http.StatusBadRequest)
			return
		}

		// Create a context with the buyer ID
		uCtx := uContext.Context{
			Properties: map[string]string{
				"buyerID": id.String(),
			},
		}

		buyer := unleashClient.GetVariant("buyerFF", unleash.WithVariantContext(uCtx))
		b, _ := json.MarshalIndent(buyer, "", "  ")
		fmt.Fprintln(w, string(b))

		if buyer.FeatureEnabled {
			fmt.Fprintln(w, "buyer is enabled!")
		} else {
			fmt.Fprintln(w, "buyer is disabled")
		}
	})

	router.HandleFunc("/buyer/{buyerID}/merchant/{merchantID}/aprs", func(w http.ResponseWriter, r *http.Request) {
		buyerID := r.PathValue("buyerID")
		merchantID := r.PathValue("merchantID")

		// Create a context with the buyer ID
		uCtx := uContext.Context{
			Properties: map[string]string{
				"buyerID":    buyerID,
				"merchantID": merchantID,
			},
		}

		buyer := unleashClient.GetVariant("test-json-for-id", unleash.WithVariantContext(uCtx))
		b, _ := json.MarshalIndent(buyer, "", "  ")
		fmt.Fprintln(w, string(b))
		fmt.Fprintln(w, buyer.Payload.Value)

		if buyer.FeatureEnabled {
			fmt.Fprintln(w, "buyer is enabled!")
		} else {
			fmt.Fprintln(w, "buyer is disabled")
		}
	})

	// Start the API server
	log.Println("API server listening on :8080")
	http.ListenAndServe(":8080", router)
}
