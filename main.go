package main

import (
	"anker/registries"
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printHelp()
		return
	}

	if len(args) == 1 {
		query := args[0]
		// Search all registries
		for _, r := range registries.All {
			searchAndPrint(r, query)
		}
	} else if len(args) == 2 {
		provider := strings.ToLower(args[0])
		query := args[1]
		found := false
		for _, r := range registries.All {
			if r.Name == provider {
				searchAndPrint(r, query)
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("Provider '%s' not found.\n", provider)
			os.Exit(1)
		}
	} else {
		fmt.Println("Too many arguments. Use --help for usage.")
		os.Exit(1)
	}
}

func searchAndPrint(reg registries.Registry, query string) {
	fmt.Printf("Searching %s for '%s'...\n", reg.Name, query)
	results, err := reg.Search(query)
	if err != nil {
		fmt.Printf("Error searching %s: %v\n", reg.Name, err)
		return
	}
	if len(results) == 0 {
		fmt.Printf("No results found on %s.\n", reg.Name)
		return
	}
	for _, res := range results {
		fmt.Printf("[%s] %s\n  Description: %s\n  URL: %s\n\n",
			reg.Name, res.Name, res.Description, res.URL)
	}
}

func printHelp() {
	fmt.Println(`Anker - Search container registries for images

Usage:
  anker <query>               Search all registries for an image
  anker <provider> <query>    Search a specific provider (e.g., quay, dockerhub)

Examples:
  anker llama
  anker quay argocd
  anker dockerhub redis

Available Providers:`)
	for _, r := range registries.All {
		fmt.Printf("  - %s\n", r.Name)
	}
}

