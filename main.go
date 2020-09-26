package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/wietsevenema/shutdown/lib/run"
)

type App struct {
}

func port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func main() {
	app := &App{}
	http.HandleFunc("/", app.serveIndex)
	http.HandleFunc("/shutdown", app.deleteSelf)

	// Start server
	log.Println("Listening on port " + port())
	log.Fatal(http.ListenAndServe(":"+port(), nil))
}

// serveIndex returns the index.html file
func (app *App) serveIndex(
	w http.ResponseWriter, r *http.Request) {

	type IndexPage struct {
	}

	// Render page template
	tpl := template.Must(
		template.New("index.html").
			ParseFiles("web/index.html"))
	tpl.Execute(w, &IndexPage{})
}

// serveIndex returns the index.html file
func (app *App) deleteSelf(
	w http.ResponseWriter, r *http.Request) {

	err := run.DeleteSelf()
	if err != nil {
		fmt.Fprint(w, err)
	}
}
