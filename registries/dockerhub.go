package registries

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ImageResult struct {
	Name        string
	URL         string
	Description string
}

type dockerHubResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		RepoName         string `json:"repo_name"`
		ShortDescription string `json:"short_description"`
	} `json:"results"`
}

func SearchDockerHub(query string) ([]ImageResult, error) {
	query = url.QueryEscape(query)
	baseURL := fmt.Sprintf("https://hub.docker.com/v2/search/repositories/?query=%s", query)
	nextURL := baseURL

	var results []ImageResult
	pageCount := 0
	maxPages := 10 // Docker Hub limit

	for {
		resp, err := http.Get(nextURL)
		if err != nil {
			return nil, fmt.Errorf("failed to query Docker Hub: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("unexpected status code: %d\nResponse: %s", resp.StatusCode, string(bodyBytes))
		}

		var data dockerHubResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to parse Docker Hub response: %v", err)
		}

		for _, item := range data.Results {
			repoURL := fmt.Sprintf("https://hub.docker.com/r/%s", item.RepoName)

			results = append(results, ImageResult{
				Name:        item.RepoName,
				URL:         repoURL,
				Description: item.ShortDescription,
			})
		}

		// Pagination cap
		if data.Next == "" || pageCount >= maxPages {
			break
		}

		nextURL = data.Next
		pageCount++
	}

	return results, nil
}
