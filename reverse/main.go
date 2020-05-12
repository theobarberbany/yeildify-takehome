package main

import (
	"log"
	"net/http"

	reverseserver "github.com/theobarberbany/yeildify-takehome/reverse/server"
)

func main() {

	server := &reverseserver.Server{}

	http.Handle("/reverse", server)
	// Return 200 on / for health checks.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	log.Fatal(http.ListenAndServe(":8090", nil))
}
