package db

import (
	"database/sql"
	"fmt"
)

// SaveLogo saves or updates the logo in the database
func (db *DB) SaveLogo(imageData []byte, mimeType string, uploadedBy int) error {
	// Delete existing logo if present
	_, err := db.Exec("DELETE FROM logo WHERE id = 1")
	if err != nil {
		return fmt.Errorf("failed to delete existing logo: %w", err)
	}

	// Verify user exists, otherwise use NULL for uploaded_by
	var userExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", uploadedBy).Scan(&userExists)
	if err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	// Insert new logo
	if userExists {
		_, err = db.Exec(`
			INSERT INTO logo (id, image_data, mime_type, uploaded_by, uploaded_at)
			VALUES (1, ?, ?, ?, CURRENT_TIMESTAMP)
		`, imageData, mimeType, uploadedBy)
	} else {
		_, err = db.Exec(`
			INSERT INTO logo (id, image_data, mime_type, uploaded_by, uploaded_at)
			VALUES (1, ?, ?, NULL, CURRENT_TIMESTAMP)
		`, imageData, mimeType)
	}
	
	if err != nil {
		return fmt.Errorf("failed to save logo: %w", err)
	}

	return nil
}

// GetLogo retrieves the stored logo
func (db *DB) GetLogo() (*Logo, error) {
	var logo Logo
	err := db.QueryRow(`
		SELECT id, image_data, mime_type, uploaded_at, COALESCE(uploaded_by, 0)
		FROM logo
		WHERE id = 1
	`).Scan(&logo.ID, &logo.ImageData, &logo.MimeType, &logo.UploadedAt, &logo.UploadedBy)

	if err == sql.ErrNoRows {
		return nil, nil // No logo uploaded yet
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get logo: %w", err)
	}

	return &logo, nil
}

// DeleteLogo removes the stored logo
func (db *DB) DeleteLogo() error {
	result, err := db.Exec("DELETE FROM logo WHERE id = 1")
	if err != nil {
		return fmt.Errorf("failed to delete logo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no logo found to delete")
	}

	return nil
}
