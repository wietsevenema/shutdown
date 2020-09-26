package run

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"cloud.google.com/go/compute/metadata"
	"golang.org/x/oauth2/google"
)

func region() string {
	if metadata.OnGCE() {
		// Cloud Run is a regional resource,
		// the zone is reported with the suffix
		// '-1' instead of the usual '-a', '-b',
		// or '-c'.
		// Example: europe-west1-1
		zone, err := metadata.Zone()
		if err != nil {
			log.Fatal(err)
		}
		return zone[:len(zone)-2]
	}
	return ""
}

func projectID() string {
	if metadata.OnGCE() {
		projectID, err := metadata.ProjectID()
		if err != nil {
			log.Fatal(err)
		}
		return projectID
	}
	return ""
}

func DeleteSelf() error {
	name := os.Getenv("K_SERVICE")
	if name == "" {
		return fmt.Errorf("can't determine service name")
	}

	client, err := google.DefaultClient(context.Background())
	if err != nil {
		return fmt.Errorf("setting up http client: %v", err)
	}
	// Build the URL to call. There is a separate
	// API endpoint for each Cloud Run region.
	apiURL := fmt.Sprintf("https://%s-run.googleapis.com/"+
		"apis/serving.knative.dev/v1/"+
		"namespaces/%s/"+
		"services/%v",
		region(),
		projectID(),
		name)

	requestURL, _ := url.Parse(apiURL)
	req := &http.Request{
		Method: "DELETE",
		URL:    requestURL,
	}

	// Call the API and handle errors
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("calling Run API: %v", err)
	}
	// Close the response when the surrounding
	// function returns
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("service %s not found "+
			"in region %s, ",
			name,
			region(),
		)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %v from Run API", resp.StatusCode)
	}
	return nil

}
