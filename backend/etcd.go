package backend

import (
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

const (
	prefix = "/t2"
)

type etcdb struct {
	client    *etcd.Client
	namespace string
}

func newEtcd(ns string, m []string) Backend {
	return &etcdb{
		client:    etcd.NewClient(m),
		namespace: strings.Join([]string{prefix, ns}, "/"),
	}
}

func (et *etcdb) key(k string) string {
	return strings.Join([]string{et.namespace, k}, "/")
}

func (et *etcdb) Get(key string) (val string, err error) {
	response, err := et.client.Get(et.key(key), false, false)
	if err != nil {
		return
	}

	val = response.Node.Value
	return
}

func (et *etcdb) GetAll(key string) (vals []string, err error) {
	response, err := et.client.Get(et.key(key), true, false)
	if err != nil {
		return
	}
	if response.Node.Nodes == nil {
		return
	}

	for _, node := range response.Node.Nodes {
		vals = append(vals, node.Value)
	}
	return
}
