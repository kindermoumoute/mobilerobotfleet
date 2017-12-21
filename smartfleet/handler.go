package smartfleet

import (
	"context"
	"fmt"
	"net/http"

	"encoding/json"
	"io/ioutil"

	"time"

	"strconv"

	"github.com/coreos/etcd/client"
)

func (s *SmartFleet) EndPoint(w http.ResponseWriter, r *http.Request) {
	//out, err := s.cmd.runScript()
	//if err != nil {
	//	fmt.Fprintf(w, "Erreur: %s", err)
	//} else {
	//	fmt.Fprintf(w, "Programme exécuté : %s", out)
	//}
	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Erreur: %s", err)
		return
	}
	bodyString := string(bodyBytes)
	key := time.Now()
	j := Job{
		Task: bodyString,
		States: []State{
			State{
				State: JobStateQueued,
				Step:  "The job is waiting to be taken",
				Date:  key,
			},
		},
	}
	bolB, _ := json.Marshal(j)
	kAPI := client.NewKeysAPI(s.etcdClient)
	resp, err := kAPI.Set(context.Background(), "/jobs/"+strconv.Itoa(int(key.UnixNano())), string(bolB), nil)
	if err != nil {
		fmt.Fprintf(w, "\nErreur etcd: %s", err)
	} else {
		fmt.Fprintf(w, "\n%s", resp)
	}

}
