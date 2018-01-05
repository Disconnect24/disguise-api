package disguise

import (
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"fmt"
)

func SendGridHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	err := r.ParseMultipartForm(1337)
	if err != nil {
		log.Errorf(ctx, "Unable to parse form: %v", err)
	}

	for name, values := range r.MultipartForm.File {
		if name != "headers" {
			for _, value := range values {
				log.Infof(ctx, "hi: %s => %s", name, value)
			}
		}
	}

	for name, values := range r.MultipartForm.File {
		if name != "headers" {
			for _, value := range values {
				log.Infof(ctx, "hi: %s => %s", name, value)
			}
		}
	}

	fmt.Fprint(w, "haii xd")
}
