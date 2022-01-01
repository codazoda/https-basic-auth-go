package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/crypto/bcrypt"
)

type application struct {
	config struct {
		port string
		path string
		cert string
		key  string
	}
	user struct {
		name string
		hash string
	}
}

func main() {
	// Instantiate app with some application data
	app := new(application)
	// Setup the config data
	app.config.port = "8003"
	app.config.path = "./www"
	app.config.cert = "localhost.crt"
	app.config.key = "localhost.key"
	// Setup the user data
	app.user.name = "admin"
	app.user.hash = "$2a$14$A51GrX.lqVNioLKUVQSZoulowhdFq2mrFYmc/A4Um6CuUpwRREZlO" // 1234
	// Setup some routes
	http.HandleFunc("/", app.fileHandler)
	http.HandleFunc("/hello", app.helloHandler)
	http.HandleFunc("/hash", app.hashHandler)
	// Start a server
	fmt.Printf("Server started on port %s\n", app.config.port)
	if err := http.ListenAndServeTLS(":"+app.config.port, app.config.cert, app.config.key, nil); err != nil {
		fmt.Println(err)
	}
}

// A function that authenticates the user and returns true or false
func (app *application) auth(w http.ResponseWriter, r *http.Request) bool {
	// Grab the username, password, and a verification that they were formatted ok in the request
	user, pass, ok := r.BasicAuth()
	// If the header includes credentials
	if ok {
		// If the username is correct
		if user == app.user.name {
			// If the passwords match
			if bcrypt.CompareHashAndPassword([]byte(app.user.hash), []byte(pass)) == nil {
				return true
			}
		}
	}
	// Set a header reporting unauthorized and requesting credentials
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	return false
}

func (app *application) fileHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize the user and return on failure
	if !app.auth(w, r) {
		return
	}
	// Parse the url for the path
	u, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		fmt.Println("Unable to parse url.")
	}
	fmt.Println("Serving " + app.config.path + u.Path)
	// Serve a file from the default directory
	http.ServeFile(w, r, app.config.path+u.Path)
}

// A handler that authorizes the user and says hello
func (app *application) helloHandler(w http.ResponseWriter, r *http.Request) {
	// Authorize the user and return on failure
	if !app.auth(w, r) {
		return
	}
	// Say hello
	w.Header().Add("content-type", "text/html")
	io.WriteString(w, "Hello World.\n")
}

// An example of how to hash a password and store it
func (app *application) hashHandler(w http.ResponseWriter, r *http.Request) {
	pass := r.URL.Query().Get("pass")
	bytes, _ := bcrypt.GenerateFromPassword([]byte(pass), 14)
	hash := string(bytes)
	w.Header().Add("content-type", "text/html")
	io.WriteString(w, hash)
}
