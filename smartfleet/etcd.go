package smartfleet

import (
	"log"
	"time"

	"net/url"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/embed"
	"github.com/coreos/etcd/pkg/types"
)

func (s *SmartFleet) newClient() {
	var err error
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	s.EtcdClient, err = client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	//	s.Kapi = client.NewKeysAPI(s.EtcdClient)
	//	// set "/foo" key with "bar" value
	//	log.Print("Setting '/foo' key with 'bar' value")
	//	resp, err := s.Kapi.Set(context.Background(), "/foo", "bar", nil)
	//	if err != nil {
	//		log.Fatal(err)
	//	} else {
	//		// print common key info
	//		log.Printf("Set is done. Metadata is %q\n", resp)
	//	}
	//	// get "/foo" key's value
	//	log.Print("Getting '/foo' key value")
	//	resp, err = s.Kapi.Get(context.Background(), "/foo", nil)
	//	if err != nil {
	//		log.Fatal(err)
	//	} else {
	//		// print common key info
	//		log.Printf("Get is done. Metadata is %q\n", resp)
	//		// print value
	//		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	//	}
}

func (s *SmartFleet) runServer() {
	var err error
	cfg := embed.NewConfig()
	cfg.Name = s.myIP
	cfg.Dir = cfg.Name + ".etcd"
	//listenClientURL, _ := url.Parse("http://" + s.myIP + ":2379")
	//cfg.ACUrls = []url.URL{*listenClientURL}
	//cfg.LCUrls = []url.URL{*listenClientURL}
	//listenPeerURL, _ := url.Parse("http://" + s.myIP + ":2380")
	//cfg.APUrls = []url.URL{*listenPeerURL}
	//cfg.LPUrls = []url.URL{*listenPeerURL}

	u, _ := url.Parse("http://0.0.0.0:2380")
	cfg.LPUrls = []url.URL{
		*u,
	}
	u, _ = url.Parse("http://0.0.0.0:2379")
	cfg.LCUrls = []url.URL{
		*u,
	}
	u, err = url.Parse("http://" + s.myIP + ":2380")
	if err != nil {
		log.Fatal(err)
	}

	cfg.ClusterState = embed.ClusterStateFlagNew
	cfg.InitialCluster = cfg.Name + "=" + u.String()
	cfg.APUrls, _ = types.NewURLs([]string{u.String()})
	//cfg.InitialCluster = cfg.InitialClusterFromName("")
	//cfg.Dir = "default.etcd"
	//cfg.InitialCluster += ",default=" + s.sibling[0]
	cfg.ClusterState = embed.ClusterStateFlagNew
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
