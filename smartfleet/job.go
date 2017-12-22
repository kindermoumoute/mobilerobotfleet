package smartfleet

import (
	"context"
	"log"
	"time"

	"encoding/json"

	"github.com/coreos/etcd/client"
)

const (
	JobStateQueued  = "QUEUED"
	JobStateStarted = "TAKEN"
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
	if err == nil {
		log.Println("Job "+j.Task+" done! ", string(out))
	} else {
		log.Println("Job failed: ", err)
	}
	s.job = nil
}

func (s *SmartFleet) poll() {
	// polling for jobs
	e := client.NewKeysAPI(s.etcdClient)
	resp, err := e.Get(context.Background(), "/jobs/", &client.GetOptions{Recursive: true, Sort: true})
	if err != nil {
		log.Printf("No job found")
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
		log.Printf("could not unmarshal: %s", err)
		return
	}
	job = j

	state := job.States[len(job.States)-1]
	if state.State != JobStateQueued {
		job = nil
		return
	}

	job.setState(JobStateStarted, "Allocated")
	job.Owner = s.myIP
	b, err := json.Marshal(job)
	if err != nil {
		log.Printf("ERROR: %s", err)
		job = nil
		return
	}

	// Attempt to allocate this job atomically
	_, err = e.Set(context.Background(), node.Key, string(b), nil)
	if err != nil {
		log.Printf("ERROR: %s", err)
		job = nil
	}
	return
}

func (s *SmartFleet) Work() {
	pollTick := time.Tick(10 * time.Second)
	hbTick := time.Tick(1 * time.Minute)
	stayAlive := true
	for stayAlive {
		select {
		case <-hbTick:
			//log.Print("heartbeat")
		case <-pollTick:
			//log.Print("poll tick")
			if s.job == nil {
				s.poll()
			}
		}
	}
}
