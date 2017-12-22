package smartfleet

import (
	"log"
	"time"

	"net/url"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/embed"
	"github.com/coreos/etcd/pkg/types"
)

// client runs on localhost
func (s *SmartFleet) NewClient() {
	var err error
	cfg := client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	s.etcdClient, err = client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

// server listens on all interfaces
func (s *SmartFleet) runServer(etcdPort string) {
	var err error
	cfg := embed.NewConfig()
	cfg.Name = s.myIP
	cfg.Dir = cfg.Name + ".etcd"

	u, _ := url.Parse("http://0.0.0.0:" + etcdPort)
	cfg.LPUrls = []url.URL{
		*u,
	}
	u, _ = url.Parse("http://0.0.0.0:2379")
	cfg.LCUrls = []url.URL{
		*u,
	}
	u, err = url.Parse("http://" + s.myIP + ":" + etcdPort)
	if err != nil {
		log.Fatal(err)
	}
	cfg.ClusterState = embed.ClusterStateFlagNew
	cfg.InitialCluster = cfg.Name + "=" + u.String()
	cfg.APUrls, _ = types.NewURLs([]string{u.String()})
	for _, ip := range s.peer {
		cfg.InitialCluster += "," + ip + "=http://" + ip + ":" + etcdPort
	}
	s.etcd, err = embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer s.etcd.Close()
	select {
	case <-s.etcd.Server.ReadyNotify():
		log.Printf("Server is ready!")
	case <-time.After(60 * time.Second):
		s.etcd.Server.Stop() // trigger a shutdown
		log.Printf("Server took too long to start!")
	}
	log.Fatal(<-s.etcd.Err())
}
