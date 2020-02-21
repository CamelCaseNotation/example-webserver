package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/camelcasenotation/example-webserver/pkg/api"
)

func main() {
	srv := &http.Server{
		Addr:         ":5000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      api.RootRouter(),
	}
	fmt.Println("Listening on localhost:5000")
	log.Fatal(srv.ListenAndServe())
}
