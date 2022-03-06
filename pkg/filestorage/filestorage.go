// Package filestorage
// APIs to download/store/get files???
package filestorage

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
)

type Hash struct {
	Algo string
	Sum  []byte
}

func (h Hash) String() string {
	return fmt.Sprintf("%s-%x", h.Algo, h.Sum)
}

func (h Hash) IsEmpty() bool {
	return h.Algo == "" || len(h.Sum) < 1
}

func (h Hash) Equal(hh Hash) bool {
	return h.Algo == hh.Algo && bytes.Equal(h.Sum, hh.Sum)
}

// LocationReference descriptor for file storage
// TODO: figure out how to deal with fragments
type LocationReference struct {
	// hash of file
	Hash
	// Location url to download content
	Location string
	// ContentType mime type
	ContentType string
	// Download URL
	DownloadURL string
	// Name of file without extension?
	Name string
}

type FileStorage interface {
	// Store store file to the storage backend
	// optional Hash param to early detect existing file and verify read
	Store(io.Reader, Hash) (LocationReference, error)
	// Get get file referenced by the LocationReference
	Get(LocationReference) (io.ReadCloser, error)
	// List get list of files stored
	// maybe do an iterator instead of entire list
	List() ([]LocationReference, error)
	// Validate
	Validate() error
}

// StoreFromURL download from url and store on the backend
// calls FileStorage.Store() for resp.Body
func StoreFromURL(f FileStorage, u string) (LocationReference, error) {
	resp, err := http.Get(u)
	if err != nil {
		return LocationReference{}, err
	}
	defer resp.Body.Close()
	r, err := f.Store(resp.Body, Hash{})
	r.ContentType = resp.Header.Get("Content-type")
	cdp := resp.Header.Get("Content-Disposition")
	_, param, _ := mime.ParseMediaType(cdp)
	fname, ok := param["filename"]
	if ok {
		r.Name = fname
	}

	r.DownloadURL = u
	return r, err
}
