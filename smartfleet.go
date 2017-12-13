package main

import (
	"fmt"
	"net/http"
	"runtime"
)

type smartFleet struct {
	cmd terminal
}

func newSmartFleet() (smartFleet, error) {
	s := smartFleet{}
	if runtime.GOOS == "windows" {
		s.cmd = windows{}
	} else {
		s.cmd = other{}
	}
	// TODO: parser les flags en entrée

	// TODO: connection aux autres neuds RAFT

	return s, nil
}

func (s smartFleet) EndPoint(w http.ResponseWriter, r *http.Request) {
	out, err := s.cmd.runScript()
	if err != nil {
		fmt.Fprintf(w, "Erreur: %s", err)
	} else {
		fmt.Fprintf(w, "Programme exécuté : %s", out)
	}
}
