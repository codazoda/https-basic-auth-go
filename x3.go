package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
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
		pass string
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
	app.user.pass = "A6xnQhbz4Vx2HuGl4lXwZ5U2I8iziLRFnhP5eNfIRvQ=" // 1234
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
		// Base64 encode the recieved password
		sum := sha256.Sum256([]byte(pass))
		encodedPass := base64.StdEncoding.EncodeToString(sum[:])
		// If the passed credentials are correct
		if user == app.user.name && encodedPass == app.user.pass {
			return true
		}
	}
	// Sleep for some random amount of time to prevent timed attacks
	rand.Seed(time.Now().UnixNano())
	wait := rand.Intn(2500)
	fmt.Printf("Waiting %d milliseconds\n", wait)
	time.Sleep(time.Duration(wait) * time.Millisecond)
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
	sum := sha256.Sum256([]byte(pass))
	w.Header().Add("content-type", "text/html")
	io.WriteString(w, base64.StdEncoding.EncodeToString(sum[:]))
}
