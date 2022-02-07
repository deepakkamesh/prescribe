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
		vidDevice    = flag.String("video_device", "/dev/video4", "Video device for teslong")
		vidH         = flag.Int("video_height", 480, "video height")
		vidW         = flag.Int("video_width", 640, "video width")
		vidFrame     = flag.Int("video_frame_rate", 10, "video frame rate")
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
	vid := NewVideo(YUYV422, uint32(*vidW), uint32(*vidH), uint(*vidFrame), *vidDevice)

	if err := vid.StartVideoStream(); err != nil {
		log.Printf("Failed to start Teslong Camera: %v", err)
	}

	// Start HTTP service.
	http := NewServer(*mockPrint, *resPath, *httpHostPort, vid)
	if err := http.Start(); err != nil {
		log.Fatalf("HTTP start failed with %v", err)
	}
}
