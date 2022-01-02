# HTTP Basic Auth in Go

This is a template for using HTTP Basic Auth in a Go application.

## Getting Started

To get started, run these commands.

```
git clone git@github.com:codazoda/x3.git
openssl req -new -newkey rsa:2048 -nodes -keyout localhost.key -out localhost.csr
openssl x509 -req -days 365 -in localhost.csr -signkey localhost.key -out localhost.crt
go build x3
./x3
```

What those commands do is...

* Clone the repo
* Generate a self-signed certificate
* Build the binary
* Run the binary

## About

Go has a built-in `BasicAuth()` method in the `net/http` module and I use that to authenticate the user. Because password hashing is so important, I'm using the bcrypt library for hasing in my template. Encryption is important with Basic Auth so we want to serve these requests over HTTPS. I've implemented TLS for this and that's why you need to generate a certificate.

I use a struct to store the application data. It will contain the server port, the web path, and the cert and key filenames. Normally you might load the username and password from a database but I've put them in this struct to keep the code simple. The username is _admin_ and the password is _1234_ for this example.

The _auth()_ function authenticates a user. I use the bcrypt library for this because encryption is hard to get right. If the header is formatted correctly and the username is correct then we compare the hash with the password. For any requests that don't authenticate, we respond indicating that the request was unauthorized and include a header that causes the browser to prompt the user for their username and password, which it will send back with the next request.

The _fileHandler()_ function authenticates then serves static files stored in the _./www_ directory.

The _helloHandler()_ function authenticates then outputs the traditional "Hello World" text to the user.

The _hashHandler()_ function _does not_ authenticate. If you pass it a _?pass=1234_ parameter, it will print a bcrypt one-way hash for that password. You'll get a different response each time because the password is properly salted by the bcrypt library. This is an example of how you might hash the password for storage in your user database. It's also the tool I used to generate the hash that I stored in the struct.
