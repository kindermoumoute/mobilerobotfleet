package smartfleet

import (
	"log"
	"time"
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

func (s *SmartFleet) poll() {

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
