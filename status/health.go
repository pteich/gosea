package status

import (
	"fmt"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")

	status := "Ok"

	/*
		w.Write([]byte(status))
	*/

	/*
		output := strings.NewReader(status)
		io.Copy(w, output)
	*/

	_, _ = fmt.Fprintf(w, "Health: %s", status)
}
