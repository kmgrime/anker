package registries

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Quay.io Search function
func SearchQuay(query string) ([]ImageResult, error) {
	escapedQuery := url.QueryEscape(query)
	apiURL := fmt.Sprintf("https://quay.io/api/v1/find/repositories?query=%s", escapedQuery)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to query Quay.io: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	// Define response structure
	var result struct {
		Results []struct {
			Name      string `json:"name"`
			Namespace struct {
				Name string `json:"name"`
			} `json:"namespace"`
			Description string `json:"description"`
			Href        string `json:"href"`
		} `json:"results"`
	}

	// Decode the response into our struct
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Quay.io response: %v", err)
	}

	// Prepare the results for return
	var results []ImageResult
	for _, repo := range result.Results {
		fullName := fmt.Sprintf("%s/%s", repo.Namespace.Name, repo.Name)
		url := fmt.Sprintf("https://quay.io/repository%s", repo.Href) // Quay URL is built from the `href` field

		results = append(results, ImageResult{
			Name:        fullName,
			URL:         url,
			Description: repo.Description,
		})
	}

	return results, nil
}
