// Package instagram content collector
package insta

import (
	"net/http"
	"net/url"
	"time"

	"github.com/parthpower/loonabot/pkg/cookie"
	"github.com/parthpower/loonabot/pkg/insta/models"
	"github.com/parthpower/loonathearchive/pkg/collector"
)

// override from go flag
var version string

type instaCollector struct {
	httpClient *http.Client
	info       collector.Info
}

func (c *instaCollector) Fetch(u string) (collector.ContentLabels, error) {
	// TODO: check if supported URL
	uu, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	q, _ := url.ParseQuery("__a=1")
	uu.RawQuery = q.Encode()
	resp, err := c.httpClient.Do(&http.Request{
		Method: http.MethodGet,
		URL:    uu,
		Proto:  "HTTP/2",
		Header: http.Header{
			"Connection": []string{"keep-alive"},
		},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	i, err := models.GetInsta(resp.Body)
	if err != nil {
		return nil, err
	}

	m, err := i.Media()
	if err != nil {
		return nil, err
	}

	l := collector.NewContentLabels().
		AddDescription(m.Caption).
		AddDownloadURLs(m.DownloadURL).
		AddSourceURL(u).
		AddFetchTime(time.Now()).
		Add("owner", "instagram.com/"+i.Items[0].User.Username)

	return l, nil
}

func (c *instaCollector) Info() (collector.Info, error) {
	return c.info, nil
}

func NewCollector(cookies string) (collector.Collector, error) {
	instaurl, _ := url.Parse("https://instagram.com")
	jar, err := cookie.ImportFromBase64(cookies, instaurl)
	if err != nil {
		return nil, err
	}
	c := &http.Client{
		Transport: http.DefaultTransport,
		Jar:       jar,
	}
	return &instaCollector{
		httpClient: c,
		info: collector.Info{
			Name: "insta",
			SupportedURLs: []string{
				// TODO: add more
				"https://instagram.com/reels/*",
				"https://instagram.com/p/*",
			},
			Version: version,
		},
	}, nil
}
