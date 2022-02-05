package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

type Server struct {
	resPath   string
	mockPrint bool
	hostPort  string
}

// response Struct to return JSON.
type response struct {
	Err  string
	Data interface{}
}

// NewServer returns an initialized server.
func NewServer(mockPrint bool, resPath string, hostPort string) *Server {
	return &Server{
		mockPrint: mockPrint,
		resPath:   resPath,
		hostPort:  hostPort,
	}
}

// Start starts the http server.
func (s *Server) Start() error {

	// Http routers.
	http.HandleFunc("/api/status", s.status)
	http.HandleFunc("/api/genpdf", s.generatePDF)
	http.HandleFunc("/api/print", s.print)

	// TODO: Setup basic auth.

	// Serve static content from resources dir.
	fs := http.FileServer(http.Dir(s.resPath))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if checkAuth(w, r) {
			fs.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
	})

	// TODO: Setup SSL.
	return http.ListenAndServe(s.hostPort, nil)
}

// print prints the pdf file that was generated.
func (s *Server) print(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	fname := strings.TrimSpace(r.Form.Get("file"))
	fpath := fmt.Sprintf("%s%s", s.resPath, fname)

	cmd := "lp"
	arg1 := "-o print-quality=3"
	if s.mockPrint {
		cmd = "stat"
		arg1 = "-t"
	}

	out, err := exec.Command(cmd, arg1, fpath).Output()
	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Print Error: %v", err),
		})
		log.Printf("Print Error %v : %v", fname, err)
		return
	}

	writeResponse(w, &response{
		Data: string(out[:]),
	})

}

// status returns the system status including status of printer and other key system metrics.
func (s *Server) status(w http.ResponseWriter, r *http.Request) {

	// If mocking Printer, ignore and return true for status.
	if s.mockPrint {
		writeResponse(w, &response{
			Data: "ok",
		})
		return
	}

	if err := checkSystemHealth(); err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("System Error: %v", err),
		})
		log.Printf("System health check error : %v", err)
		return
	}

	writeResponse(w, &response{
		Data: "ok",
	})
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
	date := now.Format("2-Jan-2006_3:04:05_pm")
	fname := fmt.Sprintf("%s-%s.pdf", name, date)

	fpath := fmt.Sprintf("%s/prescriptions/%s", s.resPath, fname)

	if err := createPDF(name, ageSex, prescription, fpath); err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Failed to create PDF:%v", err),
		})
		log.Printf("Error creating PDF: %v", err)
		return
	}

	writeResponse(w, &response{
		Data: fmt.Sprintf("/prescriptions/%s ", fname),
	})

}

// TODO: checkSystemHealth checks the health of key system parameters.
func checkSystemHealth() error {

	// Check is printer is connected by usb using device id.
	if _, err := exec.Command("lsusb", "-d", "04a9:182b").Output(); err != nil {
		return fmt.Errorf("lsusb failed to find canon: %v", err)
	}

	return nil
}

// createPDF generates a prescription PDF and stores the pdf in the file path specified
// by fname.
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

// writeResponse writes the response json object to w. If unable to marshal
// it writes a http 500.
func writeResponse(w http.ResponseWriter, resp *response) {
	w.Header().Set("Content-Type", "application/json")
	js, e := json.Marshal(resp)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	//	log.Printf("Writing json response %s", js)
	w.Write(js)
}

// checkAuth does basic validation.
func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false
	}

	return pair[0] == "guru" && pair[1] == "Karaneeswarar"
}
