package main

import (
	"fmt"
	"harbor/registries"
	"os"
)

func main() {
	// Check if an image name is passed in
	if len(os.Args) < 2 {
		fmt.Println("Usage: harbor <image-name>")
		os.Exit(1)
	}

	imageName := os.Args[1]

	// Search Docker Hub for images
	results, err := registries.SearchDockerHub(imageName) // Handle both result and error
	if err != nil {
		fmt.Println("Error:", err) // If there's an error, print it
		os.Exit(1)
	}

	// If no results, show a message
	if len(results) == 0 {
		fmt.Println("No images found.")
		return
	}

	// Print the search results
	for _, r := range results {
		fmt.Printf("[%s] %s\n  Description: %s\n  URL: %s\n\n", "dockerhub", r.Name, r.Description, r.URL)
	}
}
