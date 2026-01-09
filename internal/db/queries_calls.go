package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// CreateCall creates a new call with apparatus and responders
func (db *DB) CreateCall(call *Call, apparatusIDs []int, responderIDs []int, responderRoles []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Auto-generate incident number if not provided
	if call.IncidentNumber == "" {
		year := call.Dispatched.Year()
		call.IncidentNumber, err = db.GetNextCallNumber(year)
		if err != nil {
			return fmt.Errorf("failed to generate call number: %w", err)
		}
	}

	// Insert call
	result, err := tx.Exec(`
		INSERT INTO calls (
			incident_number, call_type, mutual_aid, address, 
			town, location_notes, dispatched, enroute, 
			on_scene, clear, narrative, created_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, call.IncidentNumber, call.CallType, call.MutualAid,
		call.Address, call.Town, call.LocationNotes,
		call.Dispatched, call.Enroute, call.OnScene, call.Clear,
		call.Narrative, call.CreatedBy)
	
	if err != nil {
		return err
	}

	callID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	call.ID = int(callID)

	// Insert apparatus
	for _, apparatusID := range apparatusIDs {
		_, err = tx.Exec(`
			INSERT INTO call_apparatus (call_id, apparatus_id) 
			VALUES (?, ?)
		`, callID, apparatusID)
		if err != nil {
			return err
		}
	}

	// Insert responders
	for i, responderID := range responderIDs {
		role := ""
		if i < len(responderRoles) {
			role = responderRoles[i]
		}
		_, err = tx.Exec(`
			INSERT INTO call_responders (call_id, responder_id, responder_role) 
			VALUES (?, ?, ?)
		`, callID, responderID, role)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetCallByID returns a call by ID with apparatus and responders
func (db *DB) GetCallByID(id int) (*Call, []Picklist, []User, error) {
	var call Call
	err := db.QueryRow(`
		SELECT c.id, c.incident_number, c.call_type, c.mutual_aid,
		       c.address, c.town, c.location_notes,
		       c.dispatched, c.enroute, c.on_scene, c.clear,
		       c.narrative, c.created_by, c.created_at, c.updated_at
		FROM calls c WHERE c.id = ?
	`, id).Scan(&call.ID, &call.IncidentNumber, &call.CallType, &call.MutualAid,
		&call.Address, &call.Town, &call.LocationNotes,
		&call.Dispatched, &call.Enroute, &call.OnScene, &call.Clear,
		&call.Narrative, &call.CreatedBy, &call.CreatedAt, &call.UpdatedAt)

	if err != nil {
		return nil, nil, nil, err
	}

	// Get apparatus
	apparatusRows, err := db.Query(`
		SELECT p.id, p.category, p.value, p.sort_order, p.active
		FROM call_apparatus ca
		JOIN picklists p ON ca.apparatus_id = p.id
		WHERE ca.call_id = ?
		ORDER BY p.sort_order, p.value
	`, id)
	if err != nil {
		return nil, nil, nil, err
	}
	defer apparatusRows.Close()

	var apparatus []Picklist
	for apparatusRows.Next() {
		var app Picklist
		err := apparatusRows.Scan(&app.ID, &app.Category, &app.Value, &app.SortOrder, &app.Active)
		if err != nil {
			return nil, nil, nil, err
		}
		apparatus = append(apparatus, app)
	}

	// Get responders
	responderRows, err := db.Query(`
		SELECT u.id, u.first_name, u.last_name, u.position, u.active, u.created, cr.responder_role
		FROM call_responders cr
		JOIN users u ON cr.responder_id = u.id
		WHERE cr.call_id = ?
		ORDER BY u.last_name, u.first_name
	`, id)
	if err != nil {
		return nil, nil, nil, err
	}
	defer responderRows.Close()

	var responders []User
	for responderRows.Next() {
		var user User
		var responderRole sql.NullString
		err := responderRows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Position, &user.Active, &user.Created, &responderRole)
		if err != nil {
			return nil, nil, nil, err
		}
		// Note: responder_role is stored in the junction table but not used in display
		responders = append(responders, user)
	}

	return &call, apparatus, responders, nil
}

// GetRecentCalls returns recent calls with pagination
func (db *DB) GetRecentCalls(limit, offset int) ([]Call, error) {
	rows, err := db.Query(`
		SELECT id, incident_number, call_type, mutual_aid,
		       address, town, location_notes,
		       dispatched, enroute, on_scene, clear,
		       narrative, created_by, created_at, updated_at
		FROM calls
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calls []Call
	for rows.Next() {
		var call Call
		err := rows.Scan(&call.ID, &call.IncidentNumber, &call.CallType, &call.MutualAid,
			&call.Address, &call.Town, &call.LocationNotes,
			&call.Dispatched, &call.Enroute, &call.OnScene, &call.Clear,
			&call.Narrative, &call.CreatedBy, &call.CreatedAt, &call.UpdatedAt)
		if err != nil {
			return nil, err
		}
		calls = append(calls, call)
	}
	return calls, nil
}

// GetCallsByYear returns all calls for a specific year
func (db *DB) GetCallsByYear(year int) ([]Call, error) {
	// Create start and end dates for the year
	startDate := fmt.Sprintf("%d-01-01T00:00:00Z", year)
	endDate := fmt.Sprintf("%d-12-31T23:59:59Z", year)
	
	rows, err := db.Query(`
		SELECT id, incident_number, call_type, mutual_aid,
		       address, town, location_notes,
		       dispatched, enroute, on_scene, clear,
		       narrative, created_by, created_at, updated_at
		FROM calls
		WHERE dispatched >= ? AND dispatched <= ?
		ORDER BY dispatched DESC
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calls []Call
	for rows.Next() {
		var call Call
		err := rows.Scan(&call.ID, &call.IncidentNumber, &call.CallType, &call.MutualAid,
			&call.Address, &call.Town, &call.LocationNotes,
			&call.Dispatched, &call.Enroute, &call.OnScene, &call.Clear,
			&call.Narrative, &call.CreatedBy, &call.CreatedAt, &call.UpdatedAt)
		if err != nil {
			return nil, err
		}
		calls = append(calls, call)
	}
	return calls, nil
}

// GetCallYears returns all years that have calls
func (db *DB) GetCallYears() ([]int, error) {
	rows, err := db.Query(`
		SELECT DISTINCT strftime('%Y', dispatched) as year
		FROM calls
		WHERE dispatched IS NOT NULL
		ORDER BY year DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var years []int
	for rows.Next() {
		var yearStr string
		err := rows.Scan(&yearStr)
		if err != nil {
			return nil, err
		}
		var year int
		fmt.Sscanf(yearStr, "%d", &year)
		years = append(years, year)
	}
	return years, nil
}

// SearchCalls searches calls based on filters
func (db *DB) SearchCalls(filters map[string]interface{}, limit, offset int) ([]Call, error) {
	query := `
		SELECT id, incident_number, call_type, mutual_aid,
		       address, town, location_notes,
		       dispatched, enroute, on_scene, clear,
		       narrative, created_by, created_at, updated_at
		FROM calls
		WHERE 1=1
	`
	
	var args []interface{}
	argIndex := 1

	// Build dynamic WHERE clause
	if startDate, ok := filters["start_date"]; ok {
		query += " AND created_at >= ?"
		args = append(args, startDate)
		argIndex++
	}
	
	if endDate, ok := filters["end_date"]; ok {
		query += " AND created_at <= ?"
		args = append(args, endDate)
		argIndex++
	}

	if callType, ok := filters["call_type"]; ok && callType != "" {
		query += " AND call_type = ?"
		args = append(args, callType)
		argIndex++
	}

	if town, ok := filters["town"]; ok && town != "" {
		query += " AND town = ?"
		args = append(args, town)
		argIndex++
	}

	if searchText, ok := filters["search_text"]; ok && searchText != "" {
		query += " AND (address LIKE ? OR incident_number LIKE ?)"
		searchPattern := "%" + searchText.(string) + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calls []Call
	for rows.Next() {
		var call Call
		err := rows.Scan(&call.ID, &call.IncidentNumber, &call.CallType, &call.MutualAid,
			&call.Address, &call.Town, &call.LocationNotes,
			&call.Dispatched, &call.Enroute, &call.OnScene, &call.Clear,
			&call.Narrative, &call.CreatedBy, &call.CreatedAt, &call.UpdatedAt)
		if err != nil {
			return nil, err
		}
		calls = append(calls, call)
	}
	return calls, nil
}

// UpdateCall updates a call (admin only or within time limit)
func (db *DB) UpdateCall(call *Call, apparatusIDs []int, responderIDs []int, responderRoles []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update call
	_, err = tx.Exec(`
		UPDATE calls SET
			incident_number = ?, call_type = ?, mutual_aid = ?,
			address = ?, town = ?, location_notes = ?,
			dispatched = ?, enroute = ?, on_scene = ?, clear = ?,
			narrative = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, call.IncidentNumber, call.CallType, call.MutualAid,
		call.Address, call.Town, call.LocationNotes,
		call.Dispatched, call.Enroute, call.OnScene, call.Clear,
		call.Narrative, call.ID)
	
	if err != nil {
		return err
	}

	// Delete existing apparatus and responders
	_, err = tx.Exec("DELETE FROM call_apparatus WHERE call_id = ?", call.ID)
	if err != nil {
		return err
	}
	
	_, err = tx.Exec("DELETE FROM call_responders WHERE call_id = ?", call.ID)
	if err != nil {
		return err
	}

	// Re-insert apparatus
	for _, apparatusID := range apparatusIDs {
		_, err = tx.Exec(`
			INSERT INTO call_apparatus (call_id, apparatus_id) 
			VALUES (?, ?)
		`, call.ID, apparatusID)
		if err != nil {
			return err
		}
	}

	// Re-insert responders
	for i, responderID := range responderIDs {
		role := ""
		if i < len(responderRoles) {
			role = responderRoles[i]
		}
		_, err = tx.Exec(`
			INSERT INTO call_responders (call_id, responder_id, responder_role) 
			VALUES (?, ?, ?)
		`, call.ID, responderID, role)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// CanUserEditCall checks if user can edit call based on time limit and role
func (db *DB) CanUserEditCall(callID, userID int) (bool, error) {
	// Get call info and settings
	var createdBy int
	var createdAt time.Time
	err := db.QueryRow(`
		SELECT created_by, created_at FROM calls WHERE id = ?
	`, callID).Scan(&createdBy, &createdAt)
	if err != nil {
		return false, err
	}

	// Get user role
	var userRole string
	err = db.QueryRow(`
		SELECT role FROM users WHERE id = ?
	`, userID).Scan(&userRole)
	if err != nil {
		return false, err
	}

	// Admin can always edit (if setting allows)
	if userRole == "admin" {
		var adminCanEdit string
		err = db.QueryRow(`
			SELECT value FROM settings WHERE key = 'admin_can_always_edit'
		`).Scan(&adminCanEdit)
		if err == nil && strings.ToLower(adminCanEdit) == "true" {
			return true, nil
		}
	}

	// Check if user created the call and within time limit
	if createdBy != userID {
		return false, nil
	}

	var timeLimitStr string
	err = db.QueryRow(`
		SELECT value FROM settings WHERE key = 'edit_time_limit_minutes'
	`).Scan(&timeLimitStr)
	if err != nil {
		return false, err
	}

	// Parse time limit (default 30 minutes)
	var timeLimit time.Duration = 30 * time.Minute
	if timeLimitStr != "" {
		if minutes := parseTimeLimit(timeLimitStr); minutes > 0 {
			timeLimit = time.Duration(minutes) * time.Minute
		}
	}

	return time.Since(createdAt) <= timeLimit, nil
}

func parseTimeLimit(s string) int {
	// Simple parser for time limit
	var minutes int
	fmt.Sscanf(s, "%d", &minutes)
	return minutes
}

// GetNextCallNumber generates the next call number for the given year
// Format: YYYY-NNN (e.g., 2026-001, 2026-002, etc.)
func (db *DB) GetNextCallNumber(year int) (string, error) {
	prefix := fmt.Sprintf("%d-", year)
	
	var maxNumber int
	err := db.QueryRow(`
		SELECT COALESCE(MAX(CAST(SUBSTR(incident_number, 6) AS INTEGER)), 0)
		FROM calls
		WHERE incident_number LIKE ? || '%'
		AND LENGTH(incident_number) = 8
	`, prefix).Scan(&maxNumber)
	
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	
	nextNumber := maxNumber + 1
	return fmt.Sprintf("%s%03d", prefix, nextNumber), nil
}
