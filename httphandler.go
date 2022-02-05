package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

type Server struct {
	resPath string
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
func (s *Server) Start(resPath string, prescriptionsPath string, hostPort string) error {

	// Http routers.
	http.HandleFunc("/api/status", s.status)
	http.HandleFunc("/api/genpdf", s.generatePDF)

	// TODO: Setup basic auth.

	// Serve static content from resources dir.
	fs := http.FileServer(http.Dir(resPath))
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
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	name := strings.TrimSpace(r.Form.Get("name"))
	ageSex := strings.TrimSpace(r.Form.Get("age_sex"))
	prescription := strings.TrimSpace(r.Form.Get("prescription"))

	// Create name for file <name>-<timestamp>.
	loc, _ := time.LoadLocation("Asia/Kolkata") // Always print date/time in India time.
	now := time.Now().In(loc)
	date := now.Format("2-Jan-2006_3:04_pm")
	fname := fmt.Sprintf("%s-%s.pdf", name, date)

	fpath := fmt.Sprintf("./resources/prescriptions/%s", fname)

	if err := createPDF(name, ageSex, prescription, fpath); err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Failed to create PDF:%v", err),
		})
	}

	writeResponse(w, &response{
		Data: fmt.Sprintf("/prescriptions/%s ", fname),
	})

}

func createPDF(name string, ageSex string, prescription string, fname string) error {
	pdf := fpdf.New("P", "pt", "A4", "")
	pdf.SetLeftMargin(50.0)
	pdf.SetRightMargin(50.0)
	pdf.AddPage()

	// Write Header.
	pdf.SetFont("Arial", "B", 12)
	pdf.WriteAligned(0, 35, "Dr R GURUSWAMY, B.Sc., M.B.B.S., D.L.O.", "C")

	pdf.SetFont("Arial", "", 10)
	pdf.Ln(30)
	pdf.WriteAligned(0, 40, "Clinic:", "L")
	pdf.WriteAligned(0, 40, "Consulation:", "R")
	pdf.Ln(10)
	pdf.WriteAligned(0, 40, "New No. 22-C, Old No. 78-C", "L")
	pdf.WriteAligned(0, 40, "Monday - Saturday", "R")
	pdf.Ln(10)
	pdf.WriteAligned(0, 40, "Subramania Swamy Koil Street", "L")
	pdf.WriteAligned(0, 40, "9am - 12pm", "R")
	pdf.Ln(10)
	pdf.WriteAligned(0, 40, "Saidapet, Chennai - 600 015", "L")
	pdf.WriteAligned(0, 40, "9pm - 10pm", "R")
	pdf.Ln(10)
	pdf.WriteAligned(0, 40, "Appointments: 91760 80789", "L")
	pdf.WriteAligned(0, 40, "Phone/WhatsApp: 917xxxxxx", "R")
	pdf.Ln(10)
	pdf.WriteAligned(0, 40, "dr.guruswamy@gmail.com", "R")

	pdf.Line(50, 150, 500, 150)

	// Write Main Body.
	loc, _ := time.LoadLocation("Asia/Kolkata") // Always print date/time in India time.
	now := time.Now().In(loc)
	date := now.Format("2 Jan 2006  3:04 pm")

	pdf.Ln(40)
	pdf.WriteAligned(0, 40, fmt.Sprintf("%s - %s", name, ageSex), "L")
	pdf.WriteAligned(0, 40, date, "R")
	pdf.Ln(50)

	// Write prescriptions.
	lines := strings.Split(prescription, "\n")
	for _, line := range lines {
		pdf.WriteAligned(0, 40, line, "L")
		pdf.Ln(20)
	}

	// Create pdf.
	return pdf.OutputFileAndClose(fname)
}
