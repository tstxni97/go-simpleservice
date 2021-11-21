package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const port = ":9090"

func main() {
	proxy := httputil.NewSingleHostReverseProxy(
		&url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
		})
	log.Println("Proxy server starts.")
	log.Fatal(http.ListenAndServe(port, proxy))
}
