package smartfleet

import (
	"net"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"strings"

	"fmt"

	"github.com/tatsushid/go-fastping"
)

type terminal interface {
	runScript() ([]byte, error)
	getMyIP() ([]byte, error)
	getAddr(string) ([]string, error)
}

type windows struct {
}

func (w windows) runScript() ([]byte, error) {
	// "robview2_interpreter.exe" --set vitesse=100 -f test.rvw2
	return exec.Command("C:\\Program Files\\Didactic\\RobotinoView2\\bin\\robview2_interpreter.exe", "-f", "test.rvw2").Output()

}

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

func (w windows) getAddr(s string) ([]string, error) {
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

func (w other) runScript() ([]byte, error) {
	time.Sleep(3 * time.Second)
	return []byte("WORKED"), nil
}

func (w other) getMyIP() ([]byte, error) {
	return exec.Command("hostname", "-i").Output()

}

func (w other) getAddr(s string) ([]string, error) {
	var mutex = &sync.Mutex{}
	a := []string{}
	p := fastping.NewPinger()

	for i := 1; i <= 3; i++ {
		ra, err := net.ResolveIPAddr("ip4:icmp", "mobilerobotfleet_robotino_"+strconv.Itoa(i))
		if err != nil {
			return nil, err
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
