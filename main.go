package disguise

import (
	"net/http"

	"encoding/json"
	"fmt"
	"os"
)

// Config structure for `config.json`.
type Config struct {
	Domain         string
	MailInterval   int
	SendGridAPIKey string
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

	// Handle literally anything that isn't matched below
	http.HandleFunc("/", slashHandler)

	http.HandleFunc("/cgi-bin/account.cgi", Account)
	http.HandleFunc("/cgi-bin/check.cgi", checkHandler)
	http.HandleFunc("/cgi-bin/receive.cgi", Receive)
	http.HandleFunc("/cgi-bin/delete.cgi", Delete)
	http.HandleFunc("/cgi-bin/send.cgi", sendHandler)

	http.HandleFunc("/sendgrid/parse", SendGridHandler)
}

func slashHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hi, disguise-api here, identifying as ", global.Domain,
		" and asking Wiis to check in every ", global.MailInterval, " min")
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	Check(w, r, global.MailInterval)
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	Send(w, r, global)
}
