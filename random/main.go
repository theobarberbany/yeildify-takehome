package main

import (
	"log"
	"net/http"

	randomserver "github.com/theobarberbany/yeildify-takehome/random/server"
)

func main() {

	server := &randomserver.Server{}
	server.Initialise()

	http.Handle("/api", server)

	// Return 200 on / for health checks.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
