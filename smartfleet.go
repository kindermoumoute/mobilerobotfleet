package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

type smartFleet struct {
}

func newSmartFleet() (smartFleet, error) {
	// TODO: parser les flags en entrée

	// TODO: connection aux autres neuds RAFT

	return smartFleet{}, nil
}

func (s smartFleet) EndPoint(w http.ResponseWriter, r *http.Request) {
	c := exec.Command("C:\\Program Files\\Didactic\\RobotinoView2\\bin\\robview2_interpreter.exe", "-f", "test.rvw2")

	if err := c.Run(); err != nil {
		fmt.Fprintf(w, "Erreur: %s", err)
	} else {
		fmt.Fprintf(w, "Programme exécuté : %s", c.Stdout)
	}
}
