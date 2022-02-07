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
		teslongVideo = flag.String("teslong_video", "/dev/video4", "Video device for teslong")
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

	// Start Teslong Video.
	// TODO: make resolution and fps command param.
	vid := NewVideo(YUYV422, 640, 480, 10, *teslongVideo)

	if err := vid.StartVideoStream(); err != nil {
		log.Printf("Failed to start Teslong Camera: %v", err)
	}

	// Start HTTP service.
	http := NewServer(*mockPrint, *resPath, *httpHostPort, vid)
	if err := http.Start(); err != nil {
		log.Fatalf("HTTP start failed with %v", err)
	}
}
