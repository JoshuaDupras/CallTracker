package export

import (
	"fd-call-log/internal/db"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

// GenerateCallPDF generates a single call report PDF
func GenerateCallPDF(call *db.Call, apparatus []db.Picklist, responders []db.User, filename string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Fire Department Call Report")
	pdf.Ln(15)

	// Call information
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Incident Information")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Incident #:")
	pdf.Cell(140, 6, call.IncidentNumber)
	pdf.Ln(6)

	pdf.Cell(50, 6, "Call Type:")
	pdf.Cell(140, 6, call.CallType)
	pdf.Ln(6)

	pdf.Cell(50, 6, "Mutual Aid:")
	pdf.Cell(140, 6, call.MutualAid)
	pdf.Ln(10)

	// Location
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Location")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Address:")
	pdf.Cell(140, 6, call.Address)
	pdf.Ln(6)

	pdf.Cell(50, 6, "Town:")
	pdf.Cell(140, 6, call.Town)
	pdf.Ln(10)

	// Timeline
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Timeline")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Dispatched:")
	pdf.Cell(140, 6, call.Dispatched.Format("01/02/2006 15:04"))
	pdf.Ln(6)

	if call.Enroute != nil {
		pdf.Cell(50, 6, "Enroute:")
		pdf.Cell(140, 6, call.Enroute.Format("01/02/2006 15:04"))
		pdf.Ln(6)
	}

	if call.OnScene != nil {
		pdf.Cell(50, 6, "On Scene:")
		pdf.Cell(140, 6, call.OnScene.Format("01/02/2006 15:04"))
		pdf.Ln(6)
	}

	if call.Clear != nil {
		pdf.Cell(50, 6, "Clear:")
		pdf.Cell(140, 6, call.Clear.Format("01/02/2006 15:04"))
		pdf.Ln(6)
	}
	pdf.Ln(4)

	// Resources
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Resources")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(50, 6, "Apparatus:")
	apparatusText := "None"
	if len(apparatus) > 0 {
		apparatusText = ""
		for i, app := range apparatus {
			if i > 0 {
				apparatusText += ", "
			}
			apparatusText += app.Value
		}
	}
	pdf.Cell(140, 6, apparatusText)
	pdf.Ln(6)

	pdf.Cell(50, 6, "Responders:")
	respondersText := "None"
	if len(responders) > 0 {
		respondersText = ""
		for i, resp := range responders {
			if i > 0 {
				respondersText += ", "
			}
			respondersText += resp.FirstName + " " + resp.LastName
		}
	}
	pdf.Cell(140, 6, respondersText)
	pdf.Ln(10)

	// Narrative
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 8, "Narrative")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.Ln(6)
	
	// Multi-line narrative
	pdf.SetFont("Arial", "", 10)
	lines := pdf.SplitText(call.Narrative, 180)
	for _, line := range lines {
		pdf.Cell(190, 6, line)
		pdf.Ln(6)
	}

	return pdf.OutputFileAndClose(filename)
}

// GenerateCallLogPDF generates a tabular call log PDF
func GenerateCallLogPDF(calls []db.Call, filename string, startDate, endDate string) error {
	pdf := gofpdf.New("L", "mm", "A4", "") // Landscape orientation
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(277, 10, "Fire Department Call Log")
	pdf.Ln(8)
	
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(277, 6, fmt.Sprintf("Period: %s to %s", startDate, endDate))
	pdf.Ln(12)

	// Table headers
	pdf.SetFont("Arial", "B", 8)
	headers := []string{"Date", "Time", "Address", "Town", "Type", "Disposition"}
	widths := []float64{25, 15, 80, 40, 60, 50}
	
	for i, header := range headers {
		pdf.Cell(widths[i], 8, header)
	}
	pdf.Ln(8)

	// Table rows
	pdf.SetFont("Arial", "", 8)
	for _, call := range calls {
		pdf.Cell(widths[0], 6, call.CreatedAt.Format("01/02"))
		pdf.Cell(widths[1], 6, call.CreatedAt.Format("15:04"))
		pdf.Cell(widths[2], 6, call.Address)
		pdf.Cell(widths[3], 6, call.Town)
		pdf.Cell(widths[4], 6, call.CallType)
		pdf.Ln(6)
		
		// Check if we need a new page
		if pdf.GetY() > 180 {
			pdf.AddPage()
			// Re-print headers
			pdf.SetFont("Arial", "B", 8)
			for i, header := range headers {
				pdf.Cell(widths[i], 8, header)
			}
			pdf.Ln(8)
			pdf.SetFont("Arial", "", 8)
		}
	}

	return pdf.OutputFileAndClose(filename)
}

// GenerateSummaryPDF generates a summary statistics PDF
func GenerateSummaryPDF(calls []db.Call, filename string, startDate, endDate string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Fire Department Summary Report")
	pdf.Ln(8)
	
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 6, fmt.Sprintf("Period: %s to %s", startDate, endDate))
	pdf.Ln(15)

	// Total calls
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 8, fmt.Sprintf("Total Calls: %d", len(calls)))
	pdf.Ln(12)

	// Statistics by call type
	callTypes := make(map[string]int)
	towns := make(map[string]int)

	for _, call := range calls {
		callTypes[call.CallType]++
		towns[call.Town]++
	}

	// Call types section
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Calls by Type")
	pdf.Ln(10)
	
	pdf.SetFont("Arial", "", 10)
	for callType, count := range callTypes {
		pdf.Cell(100, 6, callType+":")
		pdf.Cell(90, 6, fmt.Sprintf("%d", count))
		pdf.Ln(6)
	}
	pdf.Ln(6)

	// Towns section
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 8, "Calls by Town")
	pdf.Ln(10)
	
	pdf.SetFont("Arial", "", 10)
	for town, count := range towns {
		if town == "" {
			town = "Unknown"
		}
		pdf.Cell(100, 6, town+":")
		pdf.Cell(90, 6, fmt.Sprintf("%d", count))
		pdf.Ln(6)
	}

	return pdf.OutputFileAndClose(filename)
}
