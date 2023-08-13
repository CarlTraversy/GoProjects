package github

import (
	"strings"
)

type GithubReadmeProjectsUrlParser struct {
}

func (parser *GithubReadmeProjectsUrlParser) ParseReadme(content string) (urls []string) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "- [") {
			parts := strings.SplitN(line, "](https://github.com/", 2)
			if len(parts) != 2 {
				continue
			}
			urlParts := strings.SplitN(parts[1], ")", 2)
			if len(urlParts) != 2 {
				continue
			}
			urls = append(urls, "https://github.com/"+urlParts[0])
		}
	}
	return
}
