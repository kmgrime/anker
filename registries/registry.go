package registries

type Registry struct {
	Name   string
	Search func(query string) ([]ImageResult, error)
}

type ImageResult struct {
	Name        string
	URL         string
	Description string
}

var All = []Registry{
	{
		Name:   "dockerhub",
		Search: SearchDockerHub,
	},
	{
		Name:   "quay",
		Search: SearchQuay,
	},
}

