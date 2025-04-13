package registries

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func SearchQuay(query string) ([]ImageResult, error) {
	encodedQuery := url.QueryEscape(query)
	searchURL := fmt.Sprintf("https://quay.io/api/v1/find/repositories?query=%s", encodedQuery)

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to query Quay.io: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected Quay status: %d\nBody: %s", resp.StatusCode, string(body))
	}

	var data struct {
		Repositories []struct {
			Name        string `json:"name"`
			Namespace   string `json:"namespace"`
			Description string `json:"description"`
		} `json:"repositories"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode Quay response: %v", err)
	}

	var results []ImageResult
	for _, item := range data.Repositories {
		fullName := fmt.Sprintf("%s/%s", item.Namespace, item.Name)
		repoURL := fmt.Sprintf("https://quay.io/repository/%s", fullName)

		results = append(results, ImageResult{
			Name:        fullName,
			URL:         repoURL,
			Description: item.Description,
		})
	}

	return results, nil
}
