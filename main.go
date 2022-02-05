package main

import (
	"flag"
	"log"
)

func main() {
	var (
		res          = flag.String("resources", "./resources", "resources directory")
		httpHostPort = flag.String("http_port", ":8080", "host:port number for http")
	)

	flag.Parse()
	log.Print("Starting Prescribe - Remote Prescription")

	http := NewServer()
	if err := http.Start(*res, *httpHostPort); err != nil {
		log.Fatalf("HTTP start failed with %v", err)
	}

}
