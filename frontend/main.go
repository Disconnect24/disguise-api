package frontend

import (
	"net/http"
	"fmt"
)

func init() {
	http.HandleFunc("/", hai)
}

func hai(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hi from frontend.")
}