package smartfleet

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/coreos/etcd/client"
)

func (s *SmartFleet) EndPoint(w http.ResponseWriter, r *http.Request) {

	// read body request
	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Erreur: %s", err)
		return
	}
	bodyString := string(bodyBytes)

	// create new job
	key := time.Now()
	job := Job{
		Task: bodyString,
		States: []State{
			State{
				State: JobStateQueued,
				Step:  "The job is waiting to be taken",
				Date:  key,
			},
		},
	}
	jsonJob, _ := json.Marshal(job)
	jobKey := strconv.Itoa(int(key.UnixNano()))
	api := client.NewKeysAPI(s.etcdClient)

	// push new job to the queue
	resp, err := api.Set(context.Background(), "/jobs/"+jobKey, string(jsonJob), nil)
	if err != nil {
		fmt.Fprintf(w, "\netcd client error: %s", err)
	} else {
		fmt.Fprintf(w, "\n%s", resp)
	}
}
