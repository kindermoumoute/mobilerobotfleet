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
	"net/http"
	"log"
)

var (
	port = "8182"
	fleetIPs = []string{
		"172.26.X.X",
		"172.26.X.X",
	}
)


func main() {
	smartFleet,err := newSmartFleet()
	if err != nil {
		// TODO: possibilité d'envoyer l'erreur à la supervision
		panic(err)
	}

	// exposition de l'API
	http.HandleFunc("/", smartFleet.EndPoint)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
// robotino	robotino
// root		dorp6

// TODO 1 (en cours): implémenter la fausse abstraction, avec une API complète
// TODO 2: virtualiser un environnement avec les 5 robots dans le même réseau
// TODO 3: implémenter une abstraction utilisant etcd
// TODO 4: implémenter une working queue et une running queue https://stackoverflow.com/questions/34629860/how-would-you-implement-a-working-queue-in-etcd
//
// docs :
// https://godoc.org/github.com/coreos/etcd/clientv3
// https://coreos.com/etcd/docs/latest/learning/api.html