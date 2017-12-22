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

// default ports
const (
	DefaultHTTPAddr = "80"
	DefaultRaftAddr = "443"
)

var HTTPport string
var RAFTport string
var PoolIPs string

func init() {
	flag.StringVar(&HTTPport, "haddr", DefaultHTTPAddr, "Set the HTTP bind address")
	flag.StringVar(&smartfleet.MyIP, "ip", smartfleet.DefaultIP, "Set current IP address")
	flag.StringVar(&PoolIPs, "pool", smartfleet.DefaultIP, "Set current IP address")
	flag.StringVar(&RAFTport, "raddr", DefaultRaftAddr, "Set Raft bind address")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <raft-data-path> \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	// set config and run etcd
	smartFleet, err := smartfleet.New(PoolIPs, RAFTport)
	if err != nil {
		panic(err)
	}

	// create etcd client
	smartFleet.NewClient()

	// expose the API
	http.HandleFunc("/", smartFleet.EndPoint)
	go func() {
		log.Fatal(http.ListenAndServe(":"+HTTPport, nil))
	}()

	// run the worker
	smartFleet.Work()
}
