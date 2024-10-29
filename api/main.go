package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Unleash/unleash-client-go/v4"
	uContext "github.com/Unleash/unleash-client-go/v4/context"
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

	router.HandleFunc("/buyer/{buyerID}/merchant/{merchantID}", func(w http.ResponseWriter, r *http.Request) {
		buyerID := r.PathValue("buyerID")
		merchantID := r.PathValue("merchantID")

		// Create a context with the buyer ID
		uCtx := uContext.Context{
			Properties: map[string]string{
				"buyerID":    buyerID,
				"merchantID": merchantID,
			},
		}

		demoFeatureFlag := unleashClient.GetVariant("demo-ff", unleash.WithVariantContext(uCtx))
		b, _ := json.MarshalIndent(demoFeatureFlag, "", "  ")
		fmt.Fprintln(w, string(b))
		fmt.Fprintln(w, demoFeatureFlag.Payload.Value)

		if demoFeatureFlag.FeatureEnabled {
			fmt.Fprintln(w, "demoFeatureFlag is enabled!")
		} else {
			fmt.Fprintln(w, "demoFeatureFlag is disabled")
		}
	})

	// Start the API server
	log.Println("API server listening on :8080")
	http.ListenAndServe(":8080", router)
}
