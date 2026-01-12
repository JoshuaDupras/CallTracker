package db

import (
	"os"
	"testing"
)

func setupTestDB(t *testing.T) (*DB, func()) {
	// Create temporary database file
	tmpFile := "test_logo.db"
	
	// Initialize database
	db, err := InitDB(tmpFile)
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	
	// Create a test admin user
	err = db.CreateUser("Test", "Admin", "Chief", "EMT", "1234", true)
	if err != nil {
		t.Logf("Note: Could not create test user (may already exist): %v", err)
	}
	
	// Return cleanup function
	cleanup := func() {
		db.Close()
		os.Remove(tmpFile)
	}
	
	return db, cleanup
}

func TestSaveLogo(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	
	// Test data - simple PNG header
	testImageData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	testMimeType := "image/png"
	testUserID := 1
	
	// Save logo
	err := db.SaveLogo(testImageData, testMimeType, testUserID)
	if err != nil {
		t.Fatalf("Failed to save logo: %v", err)
	}
	
	t.Log("Logo saved successfully")
}

func TestGetLogo(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	
	// Test data
	testImageData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	testMimeType := "image/png"
	testUserID := 1
	
	// Save logo first
	err := db.SaveLogo(testImageData, testMimeType, testUserID)
	if err != nil {
		t.Fatalf("Failed to save logo: %v", err)
	}
	
	// Retrieve logo
	logo, err := db.GetLogo()
	if err != nil {
		t.Fatalf("Failed to get logo: %v", err)
	}
	
	if logo == nil {
		t.Fatal("Expected logo to be returned, got nil")
	}
	
	// Verify data
	if logo.MimeType != testMimeType {
		t.Errorf("Expected mime type %s, got %s", testMimeType, logo.MimeType)
	}
	
	if len(logo.ImageData) != len(testImageData) {
		t.Errorf("Expected image data length %d, got %d", len(testImageData), len(logo.ImageData))
	}
	
	for i, b := range testImageData {
		if logo.ImageData[i] != b {
			t.Errorf("Image data mismatch at byte %d: expected %d, got %d", i, b, logo.ImageData[i])
		}
	}
	
	t.Logf("Logo retrieved successfully: ID=%d, MimeType=%s, Size=%d bytes", 
		logo.ID, logo.MimeType, len(logo.ImageData))
}

func TestDeleteLogo(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	
	// Test data
	testImageData := []byte{0x89, 0x50, 0x4E, 0x47}
	testMimeType := "image/png"
	testUserID := 1
	
	// Save logo first
	err := db.SaveLogo(testImageData, testMimeType, testUserID)
	if err != nil {
		t.Fatalf("Failed to save logo: %v", err)
	}
	
	// Delete logo
	err = db.DeleteLogo()
	if err != nil {
		t.Fatalf("Failed to delete logo: %v", err)
	}
	
	// Try to get logo - should return nil
	logo, err := db.GetLogo()
	if err != nil {
		t.Fatalf("Error getting logo after delete: %v", err)
	}
	
	if logo != nil {
		t.Error("Expected logo to be nil after deletion")
	}
	
	t.Log("Logo deleted successfully")
}

func TestSaveLogoReplacesExisting(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	
	// First logo
	firstImageData := []byte{0x01, 0x02, 0x03}
	firstMimeType := "image/png"
	
	err := db.SaveLogo(firstImageData, firstMimeType, 1)
	if err != nil {
		t.Fatalf("Failed to save first logo: %v", err)
	}
	
	// Second logo (should replace first)
	secondImageData := []byte{0x04, 0x05, 0x06, 0x07}
	secondMimeType := "image/jpeg"
	
	err = db.SaveLogo(secondImageData, secondMimeType, 1)
	if err != nil {
		t.Fatalf("Failed to save second logo: %v", err)
	}
	
	// Get logo and verify it's the second one
	logo, err := db.GetLogo()
	if err != nil {
		t.Fatalf("Failed to get logo: %v", err)
	}
	
	if logo.MimeType != secondMimeType {
		t.Errorf("Expected mime type %s, got %s", secondMimeType, logo.MimeType)
	}
	
	if len(logo.ImageData) != len(secondImageData) {
		t.Errorf("Expected image data length %d, got %d", len(secondImageData), len(logo.ImageData))
	}
	
	t.Log("Logo replacement works correctly")
}

func TestGetLogoWhenNoneExists(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	
	// Try to get logo when none exists
	logo, err := db.GetLogo()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if logo != nil {
		t.Error("Expected nil logo when none exists")
	}
	
	t.Log("Correctly returns nil when no logo exists")
}
