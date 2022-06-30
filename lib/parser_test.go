package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	url     string
	body    string
	summary LinkSummary
}

var commonCases = []testCase{
	{
		"http://host.com",
		"<a href=\"http://test.com\">",
		LinkSummary{
			PageUrl:          "http://host.com",
			InternalLinksNum: 0,
			ExternalLinksNum: 1,
		},
	},
}

func TestCountLinks(t *testing.T) {
	for _, c := range commonCases {
		summary := CountLinks(c.url, c.body)
		assert.Equal(t, c.summary, *summary)
	}
}

func TestIsExternal(t *testing.T) {
	assert.True(t, isExternal("test", "host"))
}

// Funny, used hasSuffix in hasScheme at first, test just for that.
func TestHasScheme(t *testing.T) {
	assert.True(t, hasScheme("http://test.com"))
}
