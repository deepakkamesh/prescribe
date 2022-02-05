package main

import (
	"flag"
	"log"
)

func main() {
	/*var (
		resPath           = flag.String("resources", "./resources", "resources directory")
		prescriptionsPath = flag.String("prescriptions_path", "./prescriptions", "prescriptions directory")
		httpHostPort      = flag.String("http_port", ":8080", "host:port number for http")
	)
	*/

	flag.Parse()
	log.Print("Starting Prescribe - Remote Prescription")

	/*	http := NewServer()
		if err := http.Start(*resPath, *prescriptionsPath, *httpHostPort); err != nil {
			log.Fatalf("HTTP start failed with %v", err)
		}*/
	pres := `Paracetomol
1 - 1 - 1 (three times a day)
antibition cream - apply externally
bland diet
no sun `

	createPDF("Mr. John Smith", "30/M", pres, "./resources/prescriptions/out.pdf")
}
