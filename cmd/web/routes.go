package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) routes(staticDir string) http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	standardChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standardChain.Then(mux)
}
