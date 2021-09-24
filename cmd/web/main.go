package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DTGlov/goweb.git/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	posts         *mysql.PostModel
	users         *mysql.UserModel
	templateCache map[string]*template.Template
}

func main() {

	//command line flag to hold the port
	addr := flag.String("addr", ":4000", "HTTP network address")

	//define a command line flag for the mysql dsn string
	dsn := flag.String("dsn", "add your username/password and the database here", "MYSQL data source name")

	//define a command line flag for the session secret
	secret := flag.String("secret", "add your secret key", "Secret key")

	//level logging for command line info
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime|log.Ldate|log.Lshortfile)

	//pass openDB the dsn string
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	//we also defer a call to db.Close(), so that the connection pool is closed before the main() fxn exits.
	defer db.Close()

	//Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	//Inititalize a new session manager
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	//initialize the app struct with the corresponding variables
	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		session:       session,
		posts:         &mysql.PostModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		templateCache: templateCache,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	//Initialize a http.Server struct with the requisite variable
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Server on port %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

//define a fxn that wraps sql.Open and returns a pool of sql connections

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	//validate the sql connection with the db.Ping call
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
