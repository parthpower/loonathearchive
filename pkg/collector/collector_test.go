package collector

import (
	"testing"
)

type dummyCollector struct {
}

func (dummyCollector) Fetch(n string) (ContentLabels, error) {
	return NewContentLabels().
			Add("something", n).
			AddDescription("test").
			AddDownloadURLs([]string{"https://github.com/parthpower/loonabot/releases/download/v0.1.1/loonabot-v0.1.1-darwin-amd64.tar.gz", "https://github.com/parthpower/loonabot/archive/refs/tags/v0.1.1.zip"}),
		nil
}

func (dummyCollector) Info() (Info, error) {
	return Info{
		Name:          "dummy",
		SupportedURLs: []string{"dummy"},
		Version:       "0.0.0+test",
	}, nil
}

func TestCollector(t *testing.T) {
	c := dummyCollector{}

	info, err := c.Info()
	if err != nil {
		t.Logf("failed to info: %q", err)
		t.Fail()
	}
	t.Logf("info: %q", info)

	l, err := c.Fetch("test")
	if err != nil {
		t.Logf("failed to fetch: %q", err)
		t.Fail()
	}
	// err = l.DownloadContent()
	// if err != nil {
	// 	t.Logf("failed to DownloadContent: %q", err)
	// 	t.Fail()
	// }
	t.Logf("fetched content labels: %q", l)
}
