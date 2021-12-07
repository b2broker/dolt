package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	uri, err := url.Parse(os.Getenv("HEALTHURI"))
	if err != nil || uri.Host == "" {
		uri = &url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
			Path:   "/health",
		}
	}
	res, err := http.Get(uri.String())
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode < 200 && res.StatusCode > 299 {
		log.Fatalln(res.StatusCode)
	}
	log.Println("OK")
}
