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
		enVideo      = flag.Bool("enable_video", false, "Enables internal video")
		vidDevice    = flag.String("video_device", "/dev/teslongcam", "Video device for teslong")
		vidH         = flag.Int("video_height", 480, "video height")
		vidW         = flag.Int("video_width", 640, "video width")
		vidFrame     = flag.Int("video_frame_rate", 10, "video frame rate")
		altVideo     = flag.String("alt_video_url", "http://192.168.0.108:8888", "Alternate Video url")
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

	var vid *Video
	// Start Teslong Video.
	if *enVideo {
		vid = NewVideo(YUYV422, uint32(*vidW), uint32(*vidH), uint(*vidFrame), *vidDevice)
	}

	// Start HTTP service.
	http := NewServer(*mockPrint, *resPath, *httpHostPort, vid, *altVideo)
	if err := http.Start(); err != nil {
		log.Fatalf("HTTP start failed with %v", err)
	}
}
