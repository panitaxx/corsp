package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/rs/cors"
)

func main() {
	targeturl := os.Getenv("URL")
	if targeturl == "" {
		log.Fatal("URL not set, please set URL environment variable")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	debug := os.Getenv("DEBUG") != ""

	target, err := url.Parse(targeturl)
	if err != nil {
		log.Fatalf("Error parsing URL: %s", err.Error())
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(target)
			pr.Out.Host = target.Host
			if debug {
				log.Printf("Proxying request %s %s", pr.In.Method, pr.In.URL.Path+pr.In.URL.RawQuery)
			}
		},
		ErrorLog: log.Default(),
		ModifyResponse: func(response *http.Response) error {
			log.Printf("<--%d %s %s", response.StatusCode, response.Request.Method, response.Request.URL.Path)
			return nil
		},
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	})

	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(":"+port, c.Handler(proxy)); err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}

}
