package lib

import (
	"fmt"
	"regexp"
	"strings"
)

type LinkSummary struct {
	PageUrl          string
	InternalLinksNum uint64
	ExternalLinksNum uint64
}

var hrefgexp = regexp.MustCompile("<a.*?href=\"(.*?)\"")

func CountLinks(host, body string) *LinkSummary {
	sum := &LinkSummary{PageUrl: host}
	matches := hrefgexp.FindAllStringSubmatch(body, -1)

	// *Assuming* that the link should be included in the count even if it's corrupted.
	// *Assuming* all the links will have an http or https scheme defined.
	for _, match := range matches {
		fmt.Println(match[1])
		if hasScheme(match[1]) && isExternal(host, match[1]) {
			sum.ExternalLinksNum += 1
		} else {
			sum.InternalLinksNum += 1
		}
	}
	return sum
}

func isExternal(host, url string) bool {
	// TODO: improve to check up until the path starts.
	return !strings.Contains(url, host)
}

func hasScheme(url string) bool {
	return strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https")
}
