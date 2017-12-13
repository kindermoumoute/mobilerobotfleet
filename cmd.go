package main

import (
	"os/exec"
	"time"
)

type terminal interface {
	runScript() ([]byte, error)
}

type windows struct {
}

func (w windows) runScript() ([]byte, error) {
	// "robview2_interpreter.exe" --set vitesse=100 -f test.rvw2
	return exec.Command("C:\\Program Files\\Didactic\\RobotinoView2\\bin\\robview2_interpreter.exe", "-f", "test.rvw2").Output()

}

type other struct {
}

func (w other) runScript() ([]byte, error) {
	time.Sleep(3 * time.Second)
	return []byte("WORKED"), nil
}
