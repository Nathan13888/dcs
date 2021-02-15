package daemon

import (
	"fmt"
	"net/http"
)

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG!")
}
