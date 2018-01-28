package disguise

import (
	"fmt"
	"google.golang.org/appengine"
	"net/http"
	"google.golang.org/appengine/log"
	"net/mail"
	"time"
)

func SendGridHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	// We sincerely hope someone won't attempt to send more than a 21MB image.
	// but, if they do, now they have 10mb for image and 1mb for text + etc
	// (still probably too much)
	err := r.ParseMultipartForm(11000000)
	if err != nil {
		log.Errorf(ctx, "Unable to parse form: %v", err)
	}

	// TODO: Check for attachments.

	if r.Form.Get("from") == "" || r.Form.Get("to") == "" || r.Form.Get("text") == "" {
		// something was nil
		log.Warningf(ctx, "Something happened to SendGrid... is someone else accessing?")
		return
	}

	wiiMail, err := formulateMail(r.Form.Get("from"), r.Form.Get("to"), r.Form.Get("text"), nil)
	if err != nil {
		log.Criticalf(ctx, "error formulating mail: %v", err)
		return
	}
	log.Infof(ctx, "mail given: \n%s", wiiMail)
	fmt.Fprint(w, "haii xd")
}

func formulateMail(from string, to string, body string, potentialImage []byte) (string, error) {
	date := time.Now().Format("02 Jan 2006 15:04:05 -0700")
	fromAddress, err := mail.ParseAddress(from)
	boundary := GenerateBoundary()
	if err != nil {
		return "", err
	}

	return fmt.Sprint(
		"Date: ", date, "\n",
		"From: ", fromAddress.Address, "\n",
		"To: ", to, "\n",
		`Content-Type: multipart/mixed; BOUNDARY="` + boundary + `"\n`,
		// todo: finish
	), nil
}
