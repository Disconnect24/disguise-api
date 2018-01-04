package disguise

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func Account(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain;charset=utf-8")
	// TODO: figure out actual mlid generation
	ctx := appengine.NewContext(r)

	r.ParseForm()

	wiiID := r.Form.Get("mlid")
	if !FriendCodeIsValid(wiiID) {
		fmt.Fprint(w, GenNormalErrorCode(610, "Invalid friend code."))
		return
	}

	taskKey := datastore.NewKey(ctx, "Accounts", wiiID[1:], 0, nil)

	// Generate passwd and mlchkid
	mlchkid := RandStringBytesMaskImprSrc(32)
	passwd := RandStringBytesMaskImprSrc(16)

	// Fill up with data.
	task := Accounts{
		Mlchkid: mlchkid,
		Passwd:  passwd,
	}

	// Saves the new entity.
	if _, err := datastore.Put(ctx, taskKey, &task); err != nil {
		log.Errorf(ctx, "Failed to save task: %v", err)
		fmt.Fprint(w, GenNormalErrorCode(450, "Database error."))
	} else {
		fmt.Fprint(w, fmt.Sprint("\n",
			GenNormalErrorCode(100, "Success."),
			"mlid=", wiiID, "\n",
			"passwd=", passwd, "\n",
			"mlchkid=", mlchkid, "\n"))
	}
}
