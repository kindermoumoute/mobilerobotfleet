package smartfleet

import (
	"runtime"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/embed"
)

const (
	statusAlive = "alive"
	DefaultIP   = "127.0.0.1"
)

var (
	MyIP = DefaultIP
	Pool = []string{
		DefaultIP,
	}
)

type SmartFleet struct {
	Kapi   client.KeysAPI
	Status string

	cmd        terminal
	myIP       string
	peer       []string
	etcd       *embed.Etcd
	etcdClient client.Client
	job        *Job
}

func New(poolIPs, raftPort string) (*SmartFleet, error) {
	Pool = []string{}
	for _, peer := range strings.Split(poolIPs, ",") {
		Pool = append(Pool, peer)
	}

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
	s.myIP = strings.Trim(s.myIP, "\r")

	s.peer, err = s.cmd.getPeersAddr(s.myIP)
	if err != nil {
		return s, err
	}

	go s.runServer(raftPort)
	return s, err
}
