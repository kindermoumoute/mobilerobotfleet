package smartfleet

import (
	"log"
	"time"
)

func (s *SmartFleet) Work() {
	log.Print("starting worker")
	pollTick := time.Tick(s.pollRate)
	hbTick := time.Tick(s.heartbeatRate)
	stayAlive := true
	// go func() {
	for stayAlive {
		select {
		case now := <-hbTick:
			log.Print("heartbeat")
			s.Heartbeat = now
			//s.save()
		case <-pollTick:
			log.Print("poll tick")
			//if s.currentJob == nil {
			//	s.poll()
			//}
		}
	}
}

// see poll function : https://github.com/OpsLabJPL/etcdq/blob/master/etcdq/worker.go
