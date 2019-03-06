package service

import(
	"net/http"
	"html/template"
)

type page struct {
  JavascriptPath string
}

func (a *App) handleHtmlRequest() http.HandlerFunc {
	t, _ := template.ParseFiles("templates/index.html")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Cache-Control", "no-cache")
		t.Execute(w, page{JavascriptPath: "xyz"})
  }
}
