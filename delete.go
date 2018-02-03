package disguise

import (
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"strconv"
	"google.golang.org/appengine/file"
	"cloud.google.com/go/storage"
)

// Delete handles delete requests of mail.
func Delete(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	// The original MySQL query specified to order by ascending.
	// Cloud Datastore seems to do this by default.
	r.ParseForm()

	delnum, err := strconv.Atoi(r.Form.Get("delnum"))
	if err != nil {
		fmt.Fprint(w, GenNormalErrorCode(ctx, 610, "delnum is invalid."))
		return
	}

	passwd := r.Form.Get("passwd")
	if passwd == "" || len(passwd) != 16 {
		fmt.Fprintf(w, GenNormalErrorCode(ctx, 330, "Unable to parse parameters."))
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

		query := datastore.NewQuery("Mail").
			Filter("Delivered = ", true).
			// Remove w from friend code
			Filter("RecipientID = ", mlidKey.StringID()).
			Limit(delnum)
		for mailToDelete := query.Run(ctx); ; {
			var currentMail Mail
			mailKey, err := mailToDelete.Next(&currentMail)
			if err == datastore.Done {
				break
			}
			if err != nil {
				log.Warningf(ctx, "Couldn't cycle through mail! %v", err)
				fmt.Fprintf(w, GenNormalErrorCode(ctx, 541, "Issue deleting mail from the database."))
				return
			}

			if currentMail.BucketedKey != "" {
				// This mail's contents is also stored in Cloud Storage.
				bucketName, err := file.DefaultBucketName(ctx)
				if err != nil {
					log.Errorf(ctx, "failed to get default bucket name: %v", err)
					fmt.Fprintf(w, GenNormalErrorCode(ctx, 541, "Issue deleting mail from the database."))
					return
				}

				// Get client
				client, err := storage.NewClient(ctx)
				if err != nil {
					log.Errorf(ctx, "failed to create storage client: %v", err)
					fmt.Fprintf(w, GenNormalErrorCode(ctx, 541, "Issue deleting mail from the database."))
					return
				}
				// In current bucket, under folder mail, get a writer for the generated filename
				err = client.Bucket(bucketName).Object("mail/" + currentMail.BucketedKey).Delete(ctx)
				if err != nil {
					log.Errorf(ctx, "failed to delete from cloud storage: %v", err)
					fmt.Fprintf(w, GenNormalErrorCode(ctx, 541, "Issue deleting mail from the database."))
				}
			}

			// delet this
			if datastore.Delete(ctx, mailKey) != nil {
				log.Errorf(ctx, "Couldn't delete mail from database!")
				fmt.Fprintf(w, GenNormalErrorCode(ctx, 541, "Issue deleting mail from the database."))
				return
			}
		}

		fmt.Fprint(w, GenNormalErrorCode(ctx, 100, "Success."),
			"delnum=", delnum)
		return
	}

	// Only runs if not returned from earlier.
	fmt.Fprintf(w, GenNormalErrorCode(ctx, 220, "Invalid authentication."))
	return
}
