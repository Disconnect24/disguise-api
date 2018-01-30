package disguise

import (
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
	"net/mail"
	//"time"
	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/file"
	"regexp"
)

func sendGridHandler(w http.ResponseWriter, r *http.Request, wiiMailDomain string) {
	ctx := appengine.NewContext(r)
	mailDomain, err := regexp.Compile(`w\d{16}\@` + wiiMailDomain)
	if err != nil {
		log.Criticalf(ctx, "error formulating regex: %v", err)
		return
	}

	// We sincerely hope someone won't attempt to send more than a 21MB image.
	// but, if they do, now they have 10mb for image and 1mb for text + etc
	// (still probably too much)
	err = r.ParseMultipartForm(11000000)
	if err != nil {
		log.Errorf(ctx, "Unable to parse form: %v", err)
		return
	}

	// TODO: Properly verify attachments.
	if r.Form.Get("from") == "" || r.Form.Get("to") == "" || r.Form.Get("text") == "" {
		// something was nil
		log.Warningf(ctx, "Something happened to SendGrid... is someone else accessing?")
		return
	}

	// Figure out who sent it.
	fromAddress, err := mail.ParseAddress(r.Form.Get("from"))
	if err != nil {
		log.Warningf(ctx, "given from address is invalid: %v", err)
		return
	}

	toAddress := r.Form.Get("to")
	// Validate who's being mailed.
	if !mailDomain.MatchString(toAddress) {
		log.Warningf(ctx, "to address didn't match")
		return
	}

	// We "create" a response for the Wii to use, based off attachments and multipart components.
	wiiMail, err := formulateMail(fromAddress.Address, toAddress, r.Form.Get("text"), nil)
	if err != nil {
		log.Criticalf(ctx, "error formulating mail: %v", err)
		return
	}
	// We'll also create a random UUID to identify the stored mail.
	mailFileName := uuid.New().String()

	// Now, we're going to begin storing that mail.
	// Get default bucket name from App Engine
	bucketName, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to get default bucket name: %v", err)
		return
	}

	// Get client
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "failed to create storage client: %v", err)
		return
	}
	defer client.Close()
	// In current bucket, under folder mail, get a writer for the generated filename
	fileWriter := client.Bucket(bucketName).Object("mail/" + mailFileName).NewWriter(ctx)
	defer fileWriter.Close()
	_, err = fileWriter.Write([]byte(wiiMail))
	if err != nil {
		log.Errorf(ctx, "failed to write wii mail to bucket: %v", err)
		return
	}

	mailKey := datastore.NewIncompleteKey(ctx, "Mail", nil)
	// Note in database
	mailStruct := Mail{
		SenderID:    fromAddress.Address,
		Body:        "",
		RecipientID: toAddress,
		Delivered:   false,
		BucketedKey: mailFileName,
	}
	_, err = datastore.Put(ctx, mailKey, &mailStruct)
	if err != nil {
		log.Errorf(ctx, "couldn't keep record of mail in database: %v", err)
		return
	}

	fmt.Fprint(w, "haii xd")
}

func formulateMail(from string, to string, body string, potentialImage []byte) (string, error) {
	//date := time.Now().Format("02 Jan 2006 15:04:05 -0700")
	//boundary := GenerateBoundary()
	//if err != nil {
	//	return "", err
	//}

	//return fmt.Sprint(
	//	"Date: ", date, "\n",
	//	"From: ", fromAddress.Address, "\n",
	//	"To: ", to, "\n",
	//	`Content-Type: multipart/mixed; BOUNDARY="` + boundary + `"\n`,
	//	so on
	//), nil
	return "henlo mailer", nil

}
