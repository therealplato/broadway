package store

import (
	"log"
	"os"
	"strings"

	etcdclient "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

var api etcdclient.KeysAPI

func init() {
	var err error
	cfg := etcdclient.Config{
		Endpoints: []string{os.Getenv("ETCD_HOST")},
	}
	client, err := etcdclient.New(cfg)
	if err != nil {
		log.Fatal("wrong etcd client")
	}
	api = etcdclient.NewKeysAPI(client)
}

type etcdStore struct{}

func New() Store {
	return &etcdStore{}
}

func (*etcdStore) SetValue(path, value string) error {
	_, err := api.Set(context.Background(), path, value, nil)
	return err
}

func (*etcdStore) Value(path string) string {
	resp, err := api.Get(context.Background(), path, nil)
	if err == nil && resp.Node != nil {
		return resp.Node.Value
	}
	return ""
}

func (*etcdStore) Values(path string) (values map[string]string) {
	values = map[string]string{}
	resp, err := api.Get(context.Background(), path, &etcdclient.GetOptions{Recursive: true})
	if err != nil {
		log.Println("Ignoring error getting values:" + path)
		log.Println(err)
		return values
	}
	if resp.Node != nil && len(resp.Node.Nodes) > 0 {
		for _, node := range resp.Node.Nodes {
			values[lastKeyItem(node.Key)] = node.Value
		}
	} else {
		log.Println("No values found here: " + path)
	}
	return values
}

func (*etcdStore) Delete(path string) error {
	_, err := api.Delete(context.Background(), path, &etcdclient.DeleteOptions{Recursive: true})
	return err
}

func lastKeyItem(key string) string {
	keyItems := strings.Split(key, "/")
	return keyItems[len(keyItems)-1]
}
