package localfs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/parthpower/loonathearchive/pkg/filestorage"
)

func TestLocalFS(t *testing.T) {
	root, err := ioutil.TempDir("", "")
	if err != nil {
		t.Logf("failed to create temp dir %q", err)
		t.FailNow()
	}
	s, err := NewLocalStorage(root)
	if err != nil {
		t.Logf("failed to create localstorage object")
		t.FailNow()
	}
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Logf("failed to create temp file %q", err)
		t.FailNow()
	}
	n := f.Name()
	f.WriteString("something test\naaaaaaaaaaaaaaaaa")
	f.Close()
	f, err = os.Open(n)
	if err != nil {
		t.Logf("failed to open temp file %q", err)
		t.FailNow()
	}
	defer f.Close()
	ref, err := s.Store(f, filestorage.Hash{})
	if err != nil {
		t.Logf("failed to store %q", err)
		t.FailNow()
	}
	t.Logf("ref: %q", ref)

	l, err := s.List()
	if err != nil || len(l) < 1 {
		t.Logf("failed to list %q", err)
		t.FailNow()
	}
	t.Logf("ref list: %q", l)

	r, err := s.Get(l[0])
	if err != nil {
		t.Logf("failed to get %q", err)
		t.FailNow()
	}
	defer r.Close()
	rr, err := ioutil.ReadAll(r)
	if err != nil {
		t.Logf("failed to readall %q", err)
		t.FailNow()
	}
	t.Logf("read back: %q", string(rr))
	s.Validate()
}
