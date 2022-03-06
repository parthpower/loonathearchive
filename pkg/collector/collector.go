// Pacakge collector
// Sturcts and interfaces for collectors
package collector

import (
	"fmt"
	"strings"
	"time"
)

type Info struct {
	Name          string
	SupportedURLs []string
	Version       string
}

type Collector interface {
	// Fetch get content from url
	Fetch(string) (ContentLabels, error)
	// Info metadata about the collector
	Info() (Info, error)
}

// ContentLabels info about content
type ContentLabels map[string]string

func NewContentLabels() ContentLabels {
	return ContentLabels{}
}

func (c ContentLabels) AddName(n string) ContentLabels {
	c["name"] = n
	return c
}

func (c ContentLabels) AddDescription(n string) ContentLabels {
	c["description"] = n
	return c
}

func (c ContentLabels) AddSourceURL(n string) ContentLabels {
	c["sourceurl"] = n
	return c
}

func (c ContentLabels) AddDownloadURLs(n []string) ContentLabels {
	c["downloadurls"] = strings.Join(n, "\n")
	return c
}

func (c ContentLabels) AddFetchTime(t time.Time) ContentLabels {
	c["fetch_timestamp"] = fmt.Sprintf("%d", t.Unix())
	return c
}

func (c ContentLabels) Add(k, v string) ContentLabels {
	c[k] = v
	return c
}

func (c ContentLabels) GetDownloadURLs() []string {
	urls, ok := c["downloadurls"]
	if !ok || urls == "" {
		return nil
	}
	return strings.Split(urls, "\n")
}
