package disguise

import (
	"fmt"
	"net/http"
	"strconv"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// Delete handles delete requests of mail.
func Delete(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	// The original MySQL query specified to order by ascending.
	// Cloud Datastore seems to do this by default.
	r.ParseForm()

	wiiID := r.Form.Get("mlid")
	if !FriendCodeIsValid(wiiID) {
		fmt.Fprint(w, GenNormalErrorCode(610, "Invalid friend code."))
		return
	}

	delnum, err := strconv.Atoi(r.Form.Get("delnum"))
	if err != nil {
		fmt.Fprint(w, GenNormalErrorCode(610, "delnum is invalid."))
		return
	}
	query := datastore.NewQuery("Mails").
		Filter("Delivered = ", true).
	// Remove w from friend code
		Filter("SenderID = ", wiiID[1:]).
		Limit(delnum)
	for mailToDelete := query.Run(ctx); ; {
		var currentMail Mail
		mailKey, err := mailToDelete.Next(&currentMail)
		if err == datastore.Done {
			break
		}
		if err != nil {
			log.Warningf(ctx, "Couldn't cycle through mail! %v", err)
			fmt.Fprintf(w, GenNormalErrorCode(541, "Issue deleting mail from the database."))
			return
		}

		// delet this
		if datastore.Delete(ctx, mailKey) != nil {
			log.Errorf(ctx, "Couldn't delete mail from database!")
			fmt.Fprintf(w, GenNormalErrorCode(541, "Issue deleting mail from the database."))
			return
		}
	}

	fmt.Fprint(w, GenNormalErrorCode(100, "Success."),
		"delnum=", delnum)
}
