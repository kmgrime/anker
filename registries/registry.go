package registries

type Registry struct {
	Name   string
	Search func(query string) ([]ImageResult, error)
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
