package etcdstore

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/parthpower/loonathearchive/pkg/mdstorage"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdMDStore struct {
	etcdClient *clientv3.Client
	rootKey    string
}

func NewEtcdMDStorage(config clientv3.Config, rootKey string) (mdstorage.MDStorage, error) {
	c, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	return &etcdMDStore{
		etcdClient: c,
		rootKey:    rootKey,
	}, nil
}

func (e *etcdMDStore) Store(id string, md *mdstorage.MD) error {
	j, err := json.Marshal(md)
	if err != nil {
		return err
	}
	_, err = e.etcdClient.Put(context.Background(), e.rootKey+"/"+id, string(j))
	return err
}

func (e *etcdMDStore) List() ([]string, error) {
	resp, err := e.etcdClient.Get(context.Background(), e.rootKey, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, err
	}
	ks := []string{}
	r := regexp.MustCompile("^" + e.rootKey + "/([a-zA-Z0-9=]*)$")
	for _, kv := range resp.Kvs {
		m := r.FindSubmatch(kv.Key)
		if len(m) < 2 {
			continue
		}
		ks = append(ks, string(m[1]))
	}
	return ks, nil
}

func (e *etcdMDStore) Get(id string) (*mdstorage.MD, error) {
	resp, err := e.etcdClient.Get(context.Background(), e.rootKey+"/"+id)
	if err != nil {
		return nil, err
	}
	if resp.Count < 1 {
		return nil, fmt.Errorf("nothing found at %s", e.rootKey+"/"+id)
	}
	var md mdstorage.MD
	err = json.Unmarshal(resp.Kvs[0].Value, &md)
	if err != nil {
		return nil, err
	}
	return &md, nil
}

// func (e *etcdMDStore) Store(id string, md *mdstorage.MD) error {
// 	var errs error
// 	_, err := e.etcdClient.Put(context.Background(), e.rootKey+"/"+id, "")
// 	if err != nil {
// 		errs = multierror.Append(errs, err)
// 	}
// 	for k, v := range md.ContentLabels {
// 		_, err := e.etcdClient.Put(context.Background(), e.rootKey+"/"+id+"/contentlabels/"+k, v, nil)
// 		if err != nil {
// 			errs = multierror.Append(errs, err)
// 		}
// 	}
// 	for i, loc := range md.LocationRefs {
// 		k := fmt.Sprintf("%s/%s/locationrefs/%d", e.rootKey, id, i)

// 		_, err := e.etcdClient.Put(context.Background(), k+"/location", loc.Location, nil)
// 		if err != nil {
// 			errs = multierror.Append(errs, err)
// 		}
// 		_, err = e.etcdClient.Put(context.Background(), k+"/contenttype", loc.ContentType, nil)
// 		if err != nil {
// 			errs = multierror.Append(errs, err)
// 		}
// 		_, err = e.etcdClient.Put(context.Background(), k+"/downloadurl", loc.DownloadURL, nil)
// 		if err != nil {
// 			errs = multierror.Append(errs, err)
// 		}
// 		_, err = e.etcdClient.Put(context.Background(), k+"/hash", loc.Hash.String(), nil)
// 		if err != nil {
// 			errs = multierror.Append(errs, err)
// 		}
// 	}
// 	return errs
// }

// func (e *etcdMDStore) List() ([]string, error) {
// 	resp, err := e.etcdClient.KV.Get(context.Background(), e.rootKey, clientv3.WithPrefix(), clientv3.WithKeysOnly())
// 	if err != nil {
// 		return nil, err
// 	}
// 	ks := []string{}
// 	r := regexp.MustCompile("^" + e.rootKey + "/([a-zA-Z0-9]/*)$")
// 	for _, kv := range resp.Kvs {
// 		m := r.FindSubmatch(kv.Key)
// 		if len(m) < 2 {
// 			continue
// 		}
// 		ks = append(ks, string(m[1]))
// 	}
// 	return ks, nil
// }

// func (e *etcdMDStore) Get(id string) (*mdstorage.MD, error) {
// 	r := regexp.MustCompile("^" + e.rootKey + "/" + id + "/" + "$")
// 	resp, err := e.etcdClient.KV.Get(context.Background(), e.rootKey+"/"+id, clientv3.WithPrefix())
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, kv := range resp.Kvs {
// 		kv.Key
// 	}

// 	return nil, nil
// }
