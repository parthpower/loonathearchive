// Package mdstorage
// metadata storage facility
package mdstorage

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/parthpower/loonathearchive/pkg/collector"
	"github.com/parthpower/loonathearchive/pkg/filestorage"
)

type MD struct {
	collector.ContentLabels
	LocationRefs []filestorage.LocationReference
}

type MDStorage interface {
	Store(string, *MD) error
	List() ([]string, error)
	Get(string) (*MD, error)
}

// DownloadStoreCollectedURL download content from urls and store them on storage backend
// gets download urls from ContentLabels.GetDownloadURLs()
// for each url calls filestorage.StoreFromURL()
// updates MD.LocationRefs
func (m *MD) DownloadStoreCollectedURL(f filestorage.FileStorage) error {
	if m == nil {
		return fmt.Errorf("md is empty")
	}
	var errs error
	errs = nil
	for _, u := range m.GetDownloadURLs() {
		ref, err := filestorage.StoreFromURL(f, u)

		if err != nil && err != os.ErrExist {
			errs = multierror.Append(errs, err)
			continue
		}
		m.LocationRefs = append(m.LocationRefs, ref)
	}
	return errs
}
