package frontend

import (
	"net/http"
	"fmt"
	"html/template"
	"os"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"io/ioutil"
)

var templates *template.Template

func init() {
	templateLocation := "./frontend/templates"
	templateDir, err := os.Open(templateLocation)
	if err != nil {
		panic(err)
	}

	templateDirList, err := templateDir.Readdir(-1)
	if err != nil {
		panic(err)
	}

	var templatePaths []string
	for _, templateFile := range templateDirList {
		templatePaths = append(templatePaths, fmt.Sprint(templateLocation, "/", templateFile.Name()))
	}

	templates, err = template.ParseFiles(templatePaths...)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// We only want the primary page.
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		s1 := templates.Lookup("header.tmpl")
		s1.ExecuteTemplate(w, "header", nil)
		fmt.Println()
		s2 := templates.Lookup("patch.tmpl")
		s2.ExecuteTemplate(w, "content", nil)
		fmt.Println()
		s3 := templates.Lookup("footer.tmpl")
		s3.ExecuteTemplate(w, "footer", nil)
		fmt.Println()
		s3.Execute(w, nil)
	})
	http.HandleFunc("/patch", configHandle)
}

func configHandle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)


	switch r.Method {
	case "POST":
		// todo: a u t h e n t i c a t i o n
		r.ParseForm()
		fileWriter, _, err := r.FormFile("uploaded_config")
		if err != nil {
			log.Errorf(ctx, "incorrect file: %v", err)
		}

		file, err := ioutil.ReadAll(fileWriter)
		patched, err := PatchNwcConfig(file)
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("Content-Disposition", "attachment; filename=\"nwc24msg.cfg\"")
		w.Write(patched)
		break
	case "GET":
		fmt.Fprint(w, "This page doesn't do anything by itself. Try going to the main site.")
	default:
		break
	}
}

func patchConfig(config []byte) error {
	return nil
}
