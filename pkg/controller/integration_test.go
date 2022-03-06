// Pacakge controller_test
// Integration tests
package controller_test

import (
	"encoding/base64"
	"io/ioutil"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/parthpower/loonathearchive/pkg/collector/insta"
	"github.com/parthpower/loonathearchive/pkg/filestorage/localfs"
	"github.com/parthpower/loonathearchive/pkg/mdstorage"
	"github.com/parthpower/loonathearchive/pkg/mdstorage/etcdstore"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var InstaCookie string

func setup() *exec.Cmd {
	os.RemoveAll("data")
	os.RemoveAll("default.etcd")
	c := exec.Command("etcd")
	c.Stdout = os.Stdout
	c.Start()
	InstaCookie = os.Getenv("INSTA_COOKIE")
	return c
}
func TestMain(m *testing.M) {
	ch := setup()
	code := m.Run()
	ch.Process.Kill()
	os.Exit(code)
}

func getLocaletcdconfig() (clientv3.Config, string) {
	return clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}, "/yeetyeet"
}

func TestInstaLocalFsEtcd(t *testing.T) {
	if InstaCookie == "" {
		t.SkipNow()
	}
	instaCollector, err := insta.NewCollector(InstaCookie)
	if err != nil {
		t.Log("failed to create insta collector: ", err)
		t.FailNow()
	}

	localFilestore, err := localfs.NewLocalStorage("data")
	if err != nil {
		t.Log("failed to create localfs filestorage: ", err)
		t.FailNow()
	}

	etcdmdstore, err := etcdstore.NewEtcdMDStorage(getLocaletcdconfig())
	if err != nil {
		t.Log("failed to create metastorage: ", err)
		t.FailNow()
	}

	testURLs := []string{
		"https://www.instagram.com/p/CaUGfTVJtF5/",
		"https://www.instagram.com/p/CaJvPyFprzo/",
		"https://www.instagram.com/p/CZ_XFQkpIe7/",
		"https://www.instagram.com/p/CZ3Zn3eP8Nv/",
	}
	// store
	for _, u := range testURLs {
		t.Log("fetching: ", u)
		c, err := instaCollector.Fetch(u)
		if err != nil {
			t.Log("failed to fetch: ", u, " error: ", err)
			t.Fail()
			continue
		}
		md := &mdstorage.MD{}
		md.ContentLabels = c
		t.Log("DownloadStoreCollectedURL: ", u)
		err = md.DownloadStoreCollectedURL(localFilestore)
		if err != nil {
			t.Log("failed to DownloadStoreCollectedURL: ", u, " error: ", err)
			t.Fail()
			continue
		}
		t.Log("etcdmdstore store: ", genId(u), " md: ", md)
		err = etcdmdstore.Store(genId(u), md)
		if err != nil {
			t.Log("failed to etcdmdstore.Store: ", u, " error: ", err)
			t.Fail()
			continue
		}
	}
	// list
	fsstoreList, err := localFilestore.List()
	if err != nil {
		t.Log("failed to localFilestore.List(): ", err)
		t.Fail()
	}
	if len(fsstoreList) < len(testURLs) {
		t.Log("fsstoreList doesn't have enough files! expected at least ", testURLs, " got ", len(fsstoreList))
		t.Fail()
	}
	t.Log("fsstoreList: ", fsstoreList)

	mdlist, err := etcdmdstore.List()
	if err != nil {
		t.Log("failed to  etcdmdstore.List().List(): ", err)
		t.Fail()
	}
	if len(mdlist) != len(testURLs) {
		t.Log("etcdmdstore doesn't have enough md! expected ", testURLs, " got ", mdlist)
		t.Fail()
	}
	t.Log("mdlist: ", mdlist)
	// get
	for _, id := range mdlist {
		md, err := etcdmdstore.Get(id)
		if err != nil {
			t.Log("failed to mdstore.Get(", id, ") err: ", err)
			t.Fail()
			continue
		}
		t.Log("md: ", md)
		for _, locref := range md.LocationRefs {
			t.Log("localFilestore.Get(", locref, ")")
			rc, err := localFilestore.Get(locref)
			if err != nil {
				t.Log("failed to get err: ", err)
				t.Fail()
				continue
			}
			defer rc.Close()
			b, err := ioutil.ReadAll(rc)
			if err != nil {
				t.Log("failed to read")
				t.Fail()
				continue
			}
			fname := locref.Name
			if fname == "" {
				fname = locref.Hash.String()
			}
			ext, err := mime.ExtensionsByType(locref.ContentType)
			if err == nil && ext != nil && len(ext) > 0 {
				fname = fname + ext[0]
			} else {
				mtype := mimetype.Lookup(locref.ContentType)
				if mtype != nil {
					fname = fname + mtype.Extension()
				}
			}
			fname = filepath.Join("data", fname)
			err = ioutil.WriteFile(fname, b, 0644)
			if err != nil {
				t.Log("failed to writefile: ", err)
				t.Fail()
			}
		}
	}

}

func genId(u string) string {
	return base64.URLEncoding.EncodeToString([]byte(u))
}
