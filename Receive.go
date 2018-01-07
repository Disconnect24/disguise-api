package disguise

import (
	"net/http"
	"strconv"

	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"math/rand"
	"time"
)

// Receive loops through stored mail and formulates a response.
// Then, if applicable, marks the mail as received.
func Receive(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Parse form.
	err := r.ParseForm()
	if err != nil {
		log.Errorf(ctx, "Error parsing form: %v", err)
	}

	maxsize, err := strconv.Atoi(r.Form.Get("maxsize"))
	if err != nil {
		log.Debugf(ctx, "maxsize given as %s. %v", r.Form.Get("maxsize"), err)
		fmt.Fprint(w, "maxsize needs to be an int.")
		return
	}

	passwd := r.Form.Get("passwd")
	if passwd == "" || len(passwd) != 16 {
		fmt.Fprintf(w, GenNormalErrorCode(330, "Unable to parse parameters."))
		return
	}

	// Check passwd
	query := datastore.NewQuery("Accounts").Filter("Passwd = ", passwd).Limit(1)
	for mlidResult := query.Run(ctx); ; {
		var currentUser Accounts
		mlidKey, err := mlidResult.Next(&currentUser)
		if err == datastore.Done {
			break
		}

		// Awesome, we're a valid user.
		// We don't need to remove the w from friend code as it's not stored that way
		mailQuery := datastore.NewQuery("Mail").
			Filter("Delivered = ", false).
			Filter("RecipientID = ", mlidKey.StringID())

		// By default, we'll assume there's no mail.
		var totalMailOutput string
		var amountOfMail = 0
		var mailSize = 0
		var wc24MimeBoundary = fmt.Sprint("BoundaryForDL", fmt.Sprint(time.Now().Format("200601021504")), "/", random(1000000, 9999999))

		// Go through returned rows and increment the size!
		for mailResult := mailQuery.Run(ctx); ; {
			var mail Mail
			mailKey, err := mailResult.Next(&mail)
			if err == datastore.Done {
				break
			}
			log.Debugf(ctx, "Mail from %s to %s", mail.SenderID, mail.RecipientID)

			individualMail := fmt.Sprint("\r\n--", wc24MimeBoundary, "\r\n")
			individualMail += "Content-Type: text/plain\r\n\r\n"
			individualMail += mail.Body

			mailSize = len(totalMailOutput + individualMail)

			// Don't add if the mail would exceed max size.
			if mailSize > maxsize {
				break
			} else {
				// Make mailSize reflect our actions.
				totalMailOutput += individualMail
				amountOfMail++

				// We're committed at this point. Mark it that way in the db.
				mail.Delivered = true
				_, err := datastore.Put(ctx, mailKey, &mail)
				if err != nil {
					log.Errorf(ctx, "Error marking mail as delivered: %v", err)
				}
			}
		}

		w.Header().Add("Content-Type", fmt.Sprint("multipart/mixed; boundary=", wc24MimeBoundary))
		fmt.Fprint(w, "--", wc24MimeBoundary, "\r\n",
			"Content-Type: text/plain\r\n\r\n",
			"This part is ignored.\r\n\r\n\r\n\n",
			"cd=100\n",
			"msg=Success.\n",
			"mailnum=", amountOfMail, "\n",
			"mailsize=", mailSize, "\n",
			"allnum=", amountOfMail, "\n",
			totalMailOutput,
			"\r\n--", wc24MimeBoundary, "--\r\n")
		return
	}

	// Only runs if not returned from earlier.
	fmt.Fprintf(w, GenNormalErrorCode(220, "Invalid authentication."))
	return
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
