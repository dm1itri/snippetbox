package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox.whendeadline.net/internal/models"
	"time"
)

type Config struct {
	addr      string
	staticDir string
}

type Application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	dsn := flag.String("dsn", "web:pass@tcp(localhost:3306)/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	app := &Application{
		infoLog:     log.New(os.Stdout, "INFO\t", log.LUTC|log.Ldate|log.Ltime),
		errorLog:    log.New(os.Stderr, "ERROR\t", log.LUTC|log.Ldate|log.Ltime|log.Lshortfile),
		formDecoder: form.NewDecoder(),
	}

	db, err := openDB(*dsn)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer db.Close()
	app.snippets = &models.SnippetModel{DB: db}
	app.users = &models.UserModel{DB: db}

	templateCache, err := newTemplateCache()
	if err != nil {
		app.errorLog.Fatal(err)
	}
	app.templateCache = templateCache

	app.sessionManager = scs.New()
	app.sessionManager.Store = mysqlstore.New(db)
	app.sessionManager.Lifetime = 12 * time.Hour
	app.sessionManager.Cookie.Secure = true

	app.infoLog.Printf("Starting server on %s\n", cfg.addr)
	server := http.Server{
		Addr:     cfg.addr,
		ErrorLog: app.errorLog,
		Handler:  app.routes(cfg.staticDir),
		TLSConfig: &tls.Config{
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
			MinVersion:       tls.VersionTLS13,
		},
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err = server.ListenAndServeTLS("tls/cert.pem", "tls/key.pem")
	app.errorLog.Fatal(err)
}
