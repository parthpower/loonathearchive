// Package localfs
// filestorage provider with local file system
package localfs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/parthpower/loonathearchive/pkg/filestorage"
)

type localStorage struct {
	root string
}

func fileExists(path string) bool {
	f, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	if f.IsDir() {
		os.Remove(path)
		return false
	}
	return true
}

func dirCreateIfDoesntExist(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewLocalStorage create new local file storage provider from root path
// if root = "" uses $PWD
// creates root dir if doesn't exist
func NewLocalStorage(root string) (filestorage.FileStorage, error) {
	err := dirCreateIfDoesntExist(root)
	if err != nil {
		return nil, err
	}
	return &localStorage{root: root}, nil
}

func (l *localStorage) Store(r io.Reader, h filestorage.Hash) (filestorage.LocationReference, error) {

	if !h.IsEmpty() && fileExists(l.getPath(h)) {
		return filestorage.LocationReference{
			Hash:     h,
			Location: filepath.Join(l.root, h.String()),
		}, nil
	}

	f, err := ioutil.TempFile(l.root, "lta")
	if err != nil {
		return filestorage.LocationReference{}, err
	}
	defer f.Close()
	defer func() { os.Remove(f.Name()) }()
	hasher := sha256.New()
	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			hasher.Write(buf[:n])
			f.Write(buf[:n])
			break
		}
		if err != nil {
			return filestorage.LocationReference{}, err
		}
		hasher.Write(buf[:n])
		_, err = f.Write(buf[:n])
		if err != nil {
			return filestorage.LocationReference{}, err
		}
	}
	sha256sum := hasher.Sum(nil)
	hash := filestorage.Hash{
		Algo: "sha256",
		Sum:  sha256sum,
	}
	// check hash with supplied hash
	if !h.IsEmpty() && !hash.Equal(h) {
		return filestorage.LocationReference{}, fmt.Errorf("file hash did not matched: supplied %s != got %s", h.String(), hash.String())
	}

	fileloc := filepath.Join(l.root, hash.String())
	if !fileExists(fileloc) {
		// move file
		name := f.Name()
		// i believe one has to close file before rename but i can be wrong
		// defer f.Close() would throw error but it's fine
		f.Close()
		err := os.Rename(name, fileloc)
		return filestorage.LocationReference{
			Hash:     hash,
			Location: fileloc,
		}, err
	}

	return filestorage.LocationReference{
		Hash:     hash,
		Location: fileloc,
	}, os.ErrExist
}

func (l *localStorage) Get(ref filestorage.LocationReference) (io.ReadCloser, error) {
	if !fileExists(ref.Location) {
		return nil, os.ErrNotExist
	}
	f, err := os.Open(ref.Location)
	return f, err
}

func (l *localStorage) List() ([]filestorage.LocationReference, error) {
	reflist := []filestorage.LocationReference{}
	filepath.Walk(l.root, func(path string, info fs.FileInfo, e error) error {
		if info.IsDir() {
			return nil
		}
		f := filepath.Base(path)
		s := strings.Split(f, "-")
		if len(s) != 2 {
			return nil
		}
		b, err := hex.DecodeString(s[1])
		if err != nil {
			return nil
		}
		reflist = append(reflist, filestorage.LocationReference{
			Hash: filestorage.Hash{
				Algo: s[0],
				Sum:  b,
			},
			Location: path,
		})
		return nil
	})
	return reflist, nil
}

func (l *localStorage) getPath(hash filestorage.Hash) string {
	return filepath.Join(l.root, hash.String())
}

// Validate check fs for hash issues
// TODO: implement Validate
func (l *localStorage) Validate() error {
	return fmt.Errorf("not implemented")
}
