package disguise

import (
	"net/http"
	"strconv"
	"fmt"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// HMAC key most likely used for `chlng`
// BUG(spotlightishere): nothing is actually done with this
var hmacKey = "ce4cf29a3d6be1c2619172b5cb298c8972d450ad"

// Check handles adding the proper interval for check.cgi along with future
// challenge solving and future mail existence checking.
// BUG(spotlightishere): Challenge solving isn't implemented whatsoever.
func Check(w http.ResponseWriter, r *http.Request, inter int) {
	ctx := appengine.NewContext(r)

	// Grab string of interval
	interval := strconv.Itoa(inter)
	// Add required headers
	w.Header().Add("Content-Type", "text/plain;charset=utf-8")
	w.Header().Add("X-Wii-Mail-Download-Span", interval)
	w.Header().Add("X-Wii-Mail-Check-Span", interval)

	// Parse form in preparation for finding mail.
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, GenNormalErrorCode(330, "Unable to parse parameters."))
		log.Errorf(ctx, "%v", err)
		return
	}

	mlchkid := r.Form.Get("mlchkid")
	if mlchkid == "" || len(mlchkid) != 32 {
		fmt.Fprintf(w, GenNormalErrorCode(330, "Unable to parse parameters."))
		return
	}

	// Check mlchkid
	query := datastore.NewQuery("Accounts").Filter("Mlchkid = ", mlchkid).Limit(1)
	for mlidResult := query.Run(ctx); ; {
		var currentUser Accounts
		_, err := mlidResult.Next(&currentUser)
		if err == datastore.Done {
			break
		}

		// Awesome, we've a valid user.
		//err = datastore.Get(ctx, userKey, currentUser)
		//if err != nil {
		//	log.Errorf(ctx, "Error loading mlid from database, despite returned from query: %v", err)
		//	fmt.Fprintf(w, GenNormalErrorCode(420, "Unable to formulate authentication statement."))
		//	return
		//}

		// We don't need to remove the w from friend code as it's not stored that way
		mailQuery := datastore.NewQuery("Mails").
			Filter("Delivered = ", false).
			Filter("SenderID = ", currentUser.Mlchkid)

		// By default, we'll assume there's no mail.
		size := 0

		// Go through returned rows and increment the size!
		for mailQuery.Run(ctx); ; {
			var mail Mail
			_, err := mlidResult.Next(&mail)
			if err == datastore.Done {
				break
			}

			size++
		}

		// mailFlag is 0 if no new mail, otherwise something random.
		var mailFlag = "0"
		if size != 0 {
			// We've more than one mail.
			mailFlag = RandStringBytesMaskImprSrc(33)
		}

		fmt.Fprint(w, GenNormalErrorCode(100, "Success."),
			"res=", hmacKey, "\n",
			"mail.flag=", mailFlag, "\n",
			"interval=", interval)
		return
	}

	// Only runs if not returned from earlier.
	fmt.Fprintf(w, GenNormalErrorCode(220, "Invalid authentication."))
	return
}
