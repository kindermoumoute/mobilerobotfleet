package smartfleet

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tatsushid/go-fastping"
)

type terminal interface {
	runScript(string) ([]byte, error)
	getMyIP() ([]byte, error)
	getPeersAddr(string) ([]string, error)
}

type windows struct {
}

var path = map[string]string{
	"A": "0",
	"B": "1",
	"C": "2",
	"D": "3",
}

func run(timeout int, command string, args ...string) string {

	// instantiate new command
	cmd := exec.Command(command, args...)

	// get pipe to standard output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "cmd.StdoutPipe() error: " + err.Error()
	}

	// start process via command
	if err := cmd.Start(); err != nil {
		return "cmd.Start() error: " + err.Error()
	}

	// setup a buffer to capture standard output
	var buf bytes.Buffer

	// create a channel to capture any errors from wait
	done := make(chan error)
	go func() {
		if _, err := buf.ReadFrom(stdout); err != nil {
			panic("buf.Read(stdout) error: " + err.Error())
		}
		done <- cmd.Wait()
	}()

	// block on select, and switch based on actions received
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			return "failed to kill: " + err.Error()
		}
		return "timeout reached, process killed"
	case err := <-done:
		if err != nil {
			close(done)
			return "process done, with error: " + err.Error()
		}
		return "process completed: " + buf.String()
	}
	return ""
}

func (w windows) runScript(s string) ([]byte, error) {
	out := run(5, "\"C:\\Program Files\\Didactic\\RobotinoView2\\bin\\robview2_interpreter.exe\"", "-f", "\"C:\\Users\\ec-lille\\robotino.rvw2\"", "--set", "path="+path[s])
	return []byte(out), nil
}

// Hacky way to get my local IP on windows
func (w windows) getMyIP() ([]byte, error) {
	b, err := exec.Command("ipconfig").Output()
	if err != nil {
		return nil, err
	}
	add := 0
	if strings.Contains(string(b), "Docker") {
		add = 1
	}
	tmp := strings.Split(string(b), "Adresse IPv4. . . . . . . . . . . . . .: ")
	tmp2 := strings.Split(tmp[1+add], "\n")
	fmt.Println("MY IP : ", tmp2[0])
	return []byte(tmp2[0]), nil
}

func (w windows) getPeersAddr(s string) ([]string, error) {
	peers := []string{}
	for _, s := range Pool {
		if s != MyIP {
			peers = append(peers, s)
		}
	}
	return peers, nil
}

type other struct {
}

func (w other) runScript(s string) ([]byte, error) {
	time.Sleep(3 * time.Second)
	return []byte("WORKED"), nil
}

func (w other) getMyIP() ([]byte, error) {
	return exec.Command("hostname", "-i").Output()

}

func (w other) getPeersAddr(s string) ([]string, error) {
	var mutex = &sync.Mutex{}
	a := []string{}
	p := fastping.NewPinger()

	for i := 1; i < 6; i++ {
		ra, err := net.ResolveIPAddr("ip4:icmp", "mobilerobotfleet_robotino_"+strconv.Itoa(i))
		if err != nil {
			break
		}
		p.AddIPAddr(ra)
	}

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		if addr.String() != s {
			mutex.Lock()
			a = append(a, addr.String())
			mutex.Unlock()
		}
	}
	return a, p.Run()
}
