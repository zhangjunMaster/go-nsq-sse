package handler

// route at "/" below in `main`.
import (
	"html/template"
	"log"
	"net/http"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	t, err := template.ParseFiles("../templates/index.html")
	if err != nil {
		log.Fatal("WTF dude, error parsing your template.")
	}
	// Render the template, writing to `w`.
	t.Execute(w, "Duder")
	log.Println("Finished HTTP request at ", r.URL.Path)
}
