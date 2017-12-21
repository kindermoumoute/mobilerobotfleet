package smartfleet

import (
	"context"
	"log"
	"time"

	"encoding/json"

	"fmt"

	"github.com/coreos/etcd/client"
)

const (
	JobStateQueued  = "QUEUED"
	JobStateStarted = "STARTED"
)

type Job struct {
	Owner  string
	Task   string
	States []State
}

type State struct {
	State string
	Step  string
	Date  time.Time
}

func (s *SmartFleet) start(j *Job) {
	out, err := s.cmd.runScript(j.Task)
	fmt.Println("TRAVAILLE TERMINE !", string(out), err)
	s.job = nil
}

func (s *SmartFleet) poll() {
	log.Print("polling for jobs")
	e := client.NewKeysAPI(s.etcdClient)
	resp, err := e.Get(context.Background(), "/jobs/", &client.GetOptions{Recursive: true, Sort: true})
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}
	var job *Job
	for _, node := range resp.Node.Nodes {
		job, err = s.tryAllocateJob(node, e)
		if err != nil {
			return
		}
		if job != nil {
			break
		}
	}
	if job != nil {
		log.Printf("starting job: %s", job.Task)
		s.job = job
		go s.start(job)
	}
}

func (job *Job) setState(jobState string, step string) {
	state := State{jobState, step, time.Now()}
	job.States = append(job.States, state)
}

func (s *SmartFleet) tryAllocateJob(node *client.Node, e client.KeysAPI) (job *Job, err error) {
	j := &Job{}
	err = json.Unmarshal([]byte(node.Value), j)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}
	job = j
	// Current state is in last position
	state := job.States[len(job.States)-1]
	if state.State != JobStateQueued {
		// Job is already allocated / in progress
		job = nil
		return
	}
	// Attempt to allocate this job atomically
	job.setState(JobStateStarted, "Allocated")
	job.Owner = s.myIP
	//job.i = path.Base(node.Key)
	b, err := json.Marshal(job)
	if err != nil {
		log.Printf("ERROR: %s", err)
		job = nil
		return
	}
	_, err = e.Set(context.Background(), node.Key, string(b), nil)
	if err != nil {
		log.Printf("ERROR: %s", err)
		job = nil
		return
	}
	return
}

func (s *SmartFleet) Work() {
	log.Print("starting worker")
	pollTick := time.Tick(s.pollRate)
	hbTick := time.Tick(s.heartbeatRate)
	stayAlive := true
	for stayAlive {
		select {
		case now := <-hbTick:
			log.Print("heartbeat")
			s.Heartbeat = now
			//s.save()
		case <-pollTick:
			log.Print("poll tick")
			if s.job == nil {
				s.poll()
			}
		}
	}
}
