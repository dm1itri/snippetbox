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

	app.infoLog.Printf("Starting server on %s\n", cfg.addr)
	server := http.Server{
		Addr:     cfg.addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(cfg.staticDir),
	}
	err := server.ListenAndServe()
	app.errorLog.Fatal(err)
}
