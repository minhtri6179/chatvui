package web

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
)

// HelloPost returns a templ component with the provided name
func HelloPost(name string) templ.Component {
	return HelloMessage(name)
}

func HelloWebHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	name := r.FormValue("name")
	component := HelloPost(name)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HelloWebHandler: %e", err)
	}
}
