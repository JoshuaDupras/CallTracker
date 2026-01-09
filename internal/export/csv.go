package export

import (
	"encoding/csv"
	"fd-call-log/internal/db"
	"fmt"
	"os"
	"strconv"
	"time"
)

// ExportCallsToCSV exports calls to CSV file
func ExportCallsToCSV(calls []db.Call, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Date", "Time", "Incident #", "Call Type", "Mutual Aid",
		"Address", "Town", "Location Notes",
		"Dispatched", "Enroute", "On Scene", "Clear",
		"Narrative", "Created By",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, call := range calls {
		record := []string{
			call.CreatedAt.Format("01/02/2006"),
			call.CreatedAt.Format("15:04"),
			call.IncidentNumber,
			call.CallType,
			call.MutualAid,
			call.Address,
			call.Town,
			call.LocationNotes,
			call.Dispatched.Format("15:04"),
			formatTimePtr(call.Enroute),
			formatTimePtr(call.OnScene),
			formatTimePtr(call.Clear),
			call.Narrative,
			strconv.Itoa(call.CreatedBy),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("15:04")
}

// BackupDatabase creates a backup of the SQLite database
func BackupDatabase(srcPath, destPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = dest.ReadFrom(src)
	return err
}

// OpenFile opens a file with the system default application
func OpenFile(filepath string) error {
	// This is a simplified version - in production you'd use
	// platform-specific code or a library like "github.com/skratchdot/open-golang/open"
	fmt.Printf("Would open file: %s\n", filepath)
	return nil
}
