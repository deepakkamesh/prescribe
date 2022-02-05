package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var (
		resPath      = flag.String("resources", "./resources", "resources directory")
		httpHostPort = flag.String("http_port", ":8080", "host:port number for http")
		mockPrint    = flag.Bool("mock_print", false, "true runs stat on file instead of printing")
		logPath      = flag.String("log_file", "./prescribe.log", "log file location")
	)

	flag.Parse()

	// Log to custom file.
	logFile, err := os.OpenFile(*logPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	// Set log out put.
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	log.Print("Starting Prescribe - Remote Prescription")

	http := NewServer(*mockPrint, *resPath, *httpHostPort)
	if err := http.Start(); err != nil {
		log.Fatalf("HTTP start failed with %v", err)
	}
}
