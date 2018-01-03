package disguise

import (
	"net/http"

	"os"
	"encoding/json"
	"fmt"
)

// Config structure for `config.json`.
type Config struct {
	Domain       string
	MailInterval int
}

var global Config

func init() {
	// Load config.
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&global)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", slashHandler)
}

func slashHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "there is to be no abuse of our christian api. proudly identifying as ", global.Domain,
		" and asking Wiis to check in every ", global.MailInterval, " min")
}
