package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Server struct {
}

// response Struct to return JSON.
type response struct {
	Err  string
	Data interface{}
}

// writeResponse writes the response json object to w. If unable to marshal
// it writes a http 500.
func writeResponse(w http.ResponseWriter, resp *response) {
	w.Header().Set("Content-Type", "application/json")
	js, e := json.Marshal(resp)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Writing json response %s", js)
	w.Write(js)
}

func NewServer() *Server {
	return &Server{}
}
func (s *Server) Start(resPath string, hostPort string) error {

	// Http routers.
	http.HandleFunc("/api/status", s.status)
	http.HandleFunc("/api/genpdf", s.generatePDF)

	// Serve static content from resources dir.
	fs := http.FileServer(http.Dir(resPath))

	// TODO: Setup basic auth.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
	// TODO: Setup SSL.
	return http.ListenAndServe(hostPort, nil)
}

// status returns the system status including status of printer and other key system metrics.
func (s *Server) status(w http.ResponseWriter, r *http.Request) {

}

// generatePDF creates the PDF file.
func (s *Server) generatePDF(w http.ResponseWriter, r *http.Request) {
	/*pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	if err := pdf.OutputFileAndClose("hello.pdf"); err != nil {
		fmt.Println(err)
	}*/
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	name := strings.TrimSpace(r.Form.Get("name"))
	ageSex := strings.TrimSpace(r.Form.Get("age_sex"))
	prescription := strings.TrimSpace(r.Form.Get("prescription"))
	fmt.Println(prescription)

	writeResponse(w, &response{
		Data: fmt.Sprintf("%s %s %s", name, ageSex, prescription),
	})

}
