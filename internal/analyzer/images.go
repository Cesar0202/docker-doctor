package analyzer

import (
	"context"
	"docker-doctor/internal/docker"
	"github.com/docker/docker/api/types/image"
)

type ImageAnalysis struct {
	Total                 int
	Dangling              int
	RecoverableSpaceBytes int64
}

func AnalyzeImages(ctx context.Context, client *docker.Client) ImageAnalysis {
	images, err := client.Cli.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		return ImageAnalysis{}
	}

	total := len(images)
	dangling := 0
	var recoverable int64

	for _, img := range images {
		if len(img.RepoTags) == 0 || (len(img.RepoTags) == 1 && img.RepoTags[0] == "<none>:<none>") {
			dangling++
			recoverable += img.Size
		}
	}

	return ImageAnalysis{
		Total:                 total,
		Dangling:              dangling,
		RecoverableSpaceBytes: recoverable,
	}
}
