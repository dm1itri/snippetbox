package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Config struct {
	addr      string
	staticDir string
}

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.Parse()

	app := &Application{
		infoLog:  log.New(os.Stdout, "INFO\t", log.LUTC|log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile),
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	app.infoLog.Printf("Starting server on %s\n", cfg.addr)
	server := http.Server{
		Addr:     cfg.addr,
		ErrorLog: app.errorLog,
		Handler:  mux}
	err := server.ListenAndServe()
	app.errorLog.Fatal(err)
}
