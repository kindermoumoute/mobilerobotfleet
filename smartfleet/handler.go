package smartfleet

import (
	"fmt"
	"net/http"
)

func (s *SmartFleet) EndPoint(w http.ResponseWriter, r *http.Request) {
	out, err := s.cmd.runScript()
	if err != nil {
		fmt.Fprintf(w, "Erreur: %s", err)
	} else {
		fmt.Fprintf(w, "Programme exécuté : %s", out)
	}
}
