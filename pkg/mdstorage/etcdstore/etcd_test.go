package etcdstore

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/parthpower/loonathearchive/pkg/collector"
	"github.com/parthpower/loonathearchive/pkg/filestorage"
	"github.com/parthpower/loonathearchive/pkg/mdstorage"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func setup() *exec.Cmd {
	c := exec.Command("etcd")
	c.Stdout = os.Stdout
	c.Start()

	return c
}

func TestMain(m *testing.M) {
	ch := setup()
	code := m.Run()
	ch.Process.Kill()
	os.Exit(code)
}

func TestEtcdStore(t *testing.T) {
	s, err := NewEtcdMDStorage(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}, "/yeehaw")
	if err != nil {
		t.Log("failed to create store: ", err)
		t.FailNow()
	}
	err = s.Store("test0", &mdstorage.MD{
		ContentLabels: collector.NewContentLabels().
			AddDescription("desc").
			AddDownloadURLs([]string{"https://github.com/parthpower/loonabot/releases/download/v0.1.1/loonabot-v0.1.1-darwin-amd64.tar.gz\nhttps://github.com/parthpower/loonabot/releases/download/v0.1.1/loonabot-v0.1.1-darwin-amd64.tar"}),
		LocationRefs: []filestorage.LocationReference{
			{
				Hash: filestorage.Hash{
					Algo: "sha256",
					Sum:  []byte{0x12, 0x00, 0x13, 0x44},
				},
				Location:    "dummy",
				ContentType: "nothing",
				DownloadURL: "something",
			},
			{
				Hash: filestorage.Hash{
					Algo: "sha256",
					Sum:  []byte{0x11, 0x03, 0x13, 0x44},
				},
				Location:    "dummy2",
				ContentType: "nothing",
				DownloadURL: "something",
			},
		},
	})
	if err != nil {
		t.Log("failed to store: ", err)
		t.FailNow()
	}
	ids, err := s.List()
	if err != nil || len(ids) < 1 {
		t.Log("failed to list: ", err)
		t.FailNow()
	}
	t.Log("ids: ", ids)
	md, err := s.Get(ids[0])
	if err != nil {
		t.Log("failed to get: ", err)
		t.FailNow()
	}
	t.Log("md: ", md)
}
