package smartfleet

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coreos/etcd/client"
)

func (s *SmartFleet) EndPoint(w http.ResponseWriter, r *http.Request) {
	out, err := s.cmd.runScript()
	if err != nil {
		fmt.Fprintf(w, "Erreur: %s", err)
	} else {
		fmt.Fprintf(w, "Programme exécuté : %s", out)
	}

	kAPI := client.NewKeysAPI(s.etcdClient)

	// create a new key /foo with the value "bar"
	_, err = kAPI.Create(context.Background(), "/foo", "bar")
	if err != nil {
		fmt.Fprintf(w, "\nErreur etcd: %s", err)
	}

	//// delete the newly created key only if the value is still "bar"
	//_, err = kAPI.Delete(context.Background(), "/foo", &client.DeleteOptions{PrevValue: "bar"})
	//if err != nil {
	//	// handle error
	//}
}
