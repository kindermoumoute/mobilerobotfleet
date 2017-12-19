//Copyright 2017 kindermoumoute
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kindermoumoute/mobilerobotfleet/smartfleet"
)

// default values
const (
	DefaultHTTPAddr = ":8080"
	DefaultRaftAddr = ":8182"
)

var HTTPport string
var PoolIPs string

//var nodeID string

//"192.168.30.106",
//"192.168.30.30",
//"192.168.30.107",

func init() {
	flag.StringVar(&HTTPport, "haddr", DefaultHTTPAddr, "Set the HTTP bind address")
	flag.StringVar(&smartfleet.MyIP, "ip", smartfleet.DefaultIP, "Set current IP address")
	flag.StringVar(&PoolIPs, "pool", smartfleet.DefaultIP, "Set current IP address")
	//flag.StringVar(&RAFTport, "raddr", DefaultRaftAddr, "Set Raft bind address")
	//flag.StringVar(&nodeID, "id", "", "Node ID")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <raft-data-path> \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	smartFleet, err := smartfleet.New(PoolIPs)
	if err != nil {
		// TODO: possibilité d'envoyer l'erreur à la supervision
		panic(err)
	}

	smartFleet.NewClient()

	// exposition de l'API
	http.HandleFunc("/", smartFleet.EndPoint)
	log.Fatal(http.ListenAndServe(HTTPport, nil))
	//smartFleet.Work()
}

// robotino	robotino
// root		dorp6

// (DONE): implémenter la fausse abstraction, avec une API complète
// (DONE): virtualiser un environnement avec les 5 robots dans le même réseau
// TODO: implémenter une abstraction utilisant etcd
// TODO: implémenter une working queue et une running queue https://stackoverflow.com/questions/34629860/how-would-you-implement-a-working-queue-in-etcd
//
// docs :
// https://godoc.org/github.com/coreos/etcd/clientv3
// https://coreos.com/etcd/docs/latest/learning/api.html
