package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

// InitDB initializes the SQLite database with schema
func InitDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	database := &DB{db}

	// Create tables
	if err := database.createTables(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// Seed default data
	if err := database.seedDefaultData(); err != nil {
		log.Printf("Warning: failed to seed default data: %v", err)
	}

	// Ensure admin user exists with PIN
	if err := database.ensureAdminExists(); err != nil {
		log.Printf("Warning: failed to ensure admin exists: %v", err)
	}

	return database, nil
}

func (db *DB) createTables() error {
	schema := `
	-- Users table
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		position TEXT NOT NULL DEFAULT 'member',
		ems_level TEXT,
		is_admin BOOLEAN NOT NULL DEFAULT 0,
		pin TEXT,
		active BOOLEAN NOT NULL DEFAULT 1,
		joined_date DATE,
		created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	-- Picklists table
	CREATE TABLE IF NOT EXISTS picklists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category TEXT NOT NULL,
		value TEXT NOT NULL,
		sort_order INTEGER NOT NULL DEFAULT 0,
		active BOOLEAN NOT NULL DEFAULT 1,
		UNIQUE(category, value)
	);

	-- Form fields configuration
	CREATE TABLE IF NOT EXISTS form_fields (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		field_name TEXT NOT NULL UNIQUE,
		label TEXT NOT NULL,
		required BOOLEAN NOT NULL DEFAULT 0,
		enabled BOOLEAN NOT NULL DEFAULT 1,
		sort_order INTEGER NOT NULL DEFAULT 0
	);

	-- Calls table
	CREATE TABLE IF NOT EXISTS calls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		incident_number TEXT,
		call_type TEXT NOT NULL,
		mutual_aid TEXT,
		address TEXT NOT NULL,
		town TEXT,
		location_notes TEXT,
		dispatched DATETIME NOT NULL,
		enroute DATETIME,
		on_scene DATETIME,
		clear DATETIME,
		narrative TEXT NOT NULL,
		created_by INTEGER NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(created_by) REFERENCES users(id)
	);

	-- Call apparatus junction table
	CREATE TABLE IF NOT EXISTS call_apparatus (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		call_id INTEGER NOT NULL,
		apparatus_id INTEGER NOT NULL,
		FOREIGN KEY(call_id) REFERENCES calls(id) ON DELETE CASCADE,
		FOREIGN KEY(apparatus_id) REFERENCES picklists(id),
		UNIQUE(call_id, apparatus_id)
	);

	-- Call responders junction table
	CREATE TABLE IF NOT EXISTS call_responders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		call_id INTEGER NOT NULL,
		responder_id INTEGER NOT NULL,
		responder_role TEXT,
		FOREIGN KEY(call_id) REFERENCES calls(id) ON DELETE CASCADE,
		FOREIGN KEY(responder_id) REFERENCES users(id),
		UNIQUE(call_id, responder_id)
	);

	-- Settings table
	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	-- Audit log table
	CREATE TABLE IF NOT EXISTS audit_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		action TEXT NOT NULL,
		table_name TEXT NOT NULL,
		record_id INTEGER,
		changes TEXT,
		timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_calls_created_at ON calls(created_at);
	CREATE INDEX IF NOT EXISTS idx_calls_call_type ON calls(call_type);
	CREATE INDEX IF NOT EXISTS idx_calls_town ON calls(town);
	CREATE INDEX IF NOT EXISTS idx_picklists_category ON picklists(category);
	CREATE INDEX IF NOT EXISTS idx_picklists_active ON picklists(active);
	`

	_, err := db.Exec(schema)
	return err
}

func (db *DB) seedDefaultData() error {
	// Check if users already exist
	var userCount int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		return err
	}
	
	// Seed default admin user if no users exist
	if userCount == 0 {
		_, err = db.Exec(`
			INSERT INTO users (first_name, last_name, position, is_admin, pin, active) VALUES 
			('Admin', 'User', 'administrator', 1, '1234', 1)
		`)
		if err != nil {
			return err
		}
	}

	// Check if picklists already exist
	var picklistCount int
	err = db.QueryRow("SELECT COUNT(*) FROM picklists").Scan(&picklistCount)
	if err != nil {
		return err
	}
	if picklistCount > 0 {
		return nil // Picklists already seeded
	}

	// Default picklist values
	picklistData := []struct {
		category   string
		values     []string
		sortOrders []int
	}{
		{"call_type", []string{"Structure Fire", "Vehicle Fire", "Grass Fire", "Medical Emergency", "Motor Vehicle Accident", "Hazmat", "Rescue", "Alarm Investigation", "Mutual Aid", "Training"}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{"mutual_aid", []string{"No", "Yes"}, []int{1, 2}},
		{"mutual_aid_agencies", []string{"Readsboro Fire Dept", "Bennington Fire Dept", "Pownal Fire Dept", "Wilmington Fire Dept", "Searsburg Fire Dept"}, []int{1, 2, 3, 4, 5}},
		{"apparatus", []string{"Engine 1", "Engine 2", "Truck 1", "Rescue 1", "Ambulance 1", "Chief", "Tanker 1"}, []int{1, 2, 3, 4, 5, 6, 7}},
		{"town", []string{"Stamford", "Readsboro", "Whitingham"}, []int{1, 2, 3}},
		{"responder_role", []string{"Driver", "Officer", "Firefighter", "EMT", "Medic", "Chief"}, []int{1, 2, 3, 4, 5, 6}},
		{"position", []string{"Chief", "Deputy Chief", "Captain", "Member", "Probationary"}, []int{1, 2, 3, 4, 5}},
		{"ems_level", []string{"None", "VEFR", "EMR", "EMT", "AEMT", "Paramedic"}, []int{1, 2, 3, 4, 5, 6}},
	}

	for _, category := range picklistData {
		for i, value := range category.values {
			_, err = db.Exec(`
				INSERT OR IGNORE INTO picklists (category, value, sort_order, active) 
				VALUES (?, ?, ?, 1)
			`, category.category, value, category.sortOrders[i])
			if err != nil {
				return err
			}
		}
	}

	// Default form fields
	formFields := []struct {
		fieldName string
		label     string
		required  bool
		enabled   bool
		sortOrder int
	}{
		{"incident_number", "Incident #", false, true, 1},
		{"call_type", "Call Type", true, true, 2},
		{"mutual_aid", "Mutual Aid", false, true, 4},
		{"address", "Address", true, true, 5},
		{"town", "Town", false, true, 7},
		{"location_notes", "Location Notes", false, true, 8},
		{"apparatus", "Apparatus", false, true, 9},
		{"responders", "Responders", false, true, 10},
		{"dispatched", "Dispatched", true, true, 11},
		{"enroute", "Enroute", false, true, 12},
		{"on_scene", "On Scene", false, true, 13},
		{"clear", "Clear", false, true, 14},
		{"narrative", "Narrative", true, true, 16},
	}

	for _, field := range formFields {
		_, err = db.Exec(`
			INSERT OR IGNORE INTO form_fields (field_name, label, required, enabled, sort_order)
			VALUES (?, ?, ?, ?, ?)
		`, field.fieldName, field.label, field.required, field.enabled, field.sortOrder)
		if err != nil {
			return err
		}
	}

	// Default settings
	_, err = db.Exec(`
		INSERT OR IGNORE INTO settings (key, value) VALUES 
		('report_dir', 'reports'),
		('auto_print_after_save', 'false'),
		('edit_time_limit_minutes', '30'),
		('admin_can_always_edit', 'true'),
		('default_date_range_days', '30')
	`)

	return err
}

// ensureAdminExists makes sure there's an admin user with a PIN
func (db *DB) ensureAdminExists() error {
	// Migrate old role-based admin to is_admin if needed
	_, err := db.Exec(`
		UPDATE users SET is_admin = 1 WHERE role = 'admin' AND is_admin = 0
	`)
	if err != nil {
		return err
	}
	
	// Check if any admin user exists
	var adminExists bool
	err = db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE is_admin = 1)
	`).Scan(&adminExists)
	
	if err != nil {
		return err
	}
	
	if !adminExists {
		// Create admin user
		_, err = db.Exec(`
			INSERT INTO users (first_name, last_name, position, is_admin, pin, active) 
			VALUES ('Admin', 'User', 'administrator', 1, '1234', 1)
		`)
		return err
	}
	
	// Check if default Admin has a PIN set
	var pin sql.NullString
	err = db.QueryRow(`
		SELECT pin FROM users WHERE first_name = 'Admin' AND last_name = 'User' AND is_admin = 1
	`).Scan(&pin)
	
	if err == sql.ErrNoRows {
		// Admin doesn't exist, but other admins do - that's ok
		return nil
	}
	if err != nil {
		return err
	}
	
	// If no PIN set, set default PIN
	if !pin.Valid || pin.String == "" {
		_, err = db.Exec(`
			UPDATE users SET pin = '1234' WHERE first_name = 'Admin' AND last_name = 'User' AND is_admin = 1
		`)
		return err
	}
	
	return nil
}
