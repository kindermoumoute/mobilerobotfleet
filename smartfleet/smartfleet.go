package smartfleet

import (
	"runtime"
	"time"

	"strings"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/embed"
)

var (
	myIP     = "192.168.30.108"
	siblings = []string{
		"192.168.30.109",
		"192.168.30.66",
	}
)

type SmartFleet struct {
	cmd     terminal
	myIP    string
	sibling []string

	EtcdClient client.Client
	Kapi       client.KeysAPI
	etcd       *embed.Etcd
}

func New() (*SmartFleet, error) {
	s := &SmartFleet{}
	if runtime.GOOS == "windows" {
		s.cmd = windows{}
	} else {
		time.Sleep(2 * time.Second)
		s.cmd = other{}
	}
	myIP, err := s.cmd.getMyIP()
	if err != nil {
		return s, err
	}
	s.myIP = strings.Trim(string(myIP), "\n")
	s.sibling, err = s.cmd.getAddr(s.myIP)
	if err != nil {
		return s, err
	}

	for i, ip := range s.sibling {
		s.sibling[i] = "http://" + ip + ":2380"
	}

	// run etcd server
	go s.runServer()

	return s, err
}
