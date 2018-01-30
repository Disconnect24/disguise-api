package disguise

import (
	"net/http"

	"encoding/json"
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

	http.HandleFunc("/cgi-bin/account.cgi", Account)
	http.HandleFunc("/cgi-bin/check.cgi", checkHandler)
	http.HandleFunc("/cgi-bin/receive.cgi", Receive)
	http.HandleFunc("/cgi-bin/delete.cgi", Delete)
	http.HandleFunc("/cgi-bin/send.cgi", sendHandler)
	http.HandleFunc("/sendgrid/parse", sendGridHandlerWrapper)
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	Check(w, r, global.MailInterval)
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	Send(w, r, global)
}

func sendGridHandlerWrapper(w http.ResponseWriter, r *http.Request) {
	sendGridHandler(w, r, global.Domain)
}
