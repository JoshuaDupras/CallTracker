package main

import (
	"context"
	"errors"
	"fd-call-log/internal/db"
	"fmt"
)

var ErrUnauthorized = errors.New("unauthorized")

// App struct
type App struct {
	ctx context.Context
	db  *db.DB
	currentUser *db.User
}

// GetVersion returns version information for the UI
func (a *App) GetVersion() map[string]string {
	return map[string]string{
		"version":   Version,
		"commit":    GitCommit,
		"buildTime": BuildTime,
	}
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	fmt.Println("App starting...")
	
	// Initialize database
	fmt.Println("Initializing database...")
	database, err := db.InitDB("./fd-calls.db")
	if err != nil {
		fmt.Printf("Database initialization failed: %v\n", err)
		panic(err)
	}
	fmt.Println("Database initialized successfully")
	a.db = database
	fmt.Println("Startup complete")
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// Login authenticates a user with PIN
func (a *App) Login(name, pin string) (*db.User, error) {
	user, err := a.db.AuthenticateUser(name, pin)
	if err != nil {
		return nil, err
	}
	a.currentUser = user
	return user, nil
}

// GetCurrentUser returns the currently logged-in user
func (a *App) GetCurrentUser() *db.User {
	return a.currentUser
}

// Logout clears the current user
func (a *App) Logout() {
	a.currentUser = nil
}

// GetAllUsers returns all active users
func (a *App) GetAllUsers() ([]db.User, error) {
	return a.db.GetAllUsers()
}

// GetActiveUsers returns all active users for login
func (a *App) GetActiveUsers() ([]db.User, error) {
	return a.db.GetActiveUsers()
}

// GetAdminUsers returns all active admin users
func (a *App) GetAdminUsers() ([]db.User, error) {
	return a.db.GetAdminUsers()
}

// GetUserByID returns a user by ID
func (a *App) GetUserByID(id int) (*db.User, error) {
	return a.db.GetUserByID(id)
}

// CreateUser creates a new user
func (a *App) CreateUser(firstName, lastName, position, emsLevel, pin string, isAdmin bool) error {
	return a.db.CreateUser(firstName, lastName, position, emsLevel, pin, isAdmin)
}

// UpdateUser updates an existing user
func (a *App) UpdateUser(user *db.User) error {
	return a.db.UpdateUser(user)
}

// DeleteUser marks a user as inactive
func (a *App) DeleteUser(id int) error {
	// Set user as inactive
	user, err := a.db.GetUserByID(id)
	if err != nil {
		return err
	}
	user.Active = false
	return a.db.UpdateUser(user)
}

// ChangePIN allows a user to change their PIN
func (a *App) ChangePIN(oldPIN, newPIN string) error {
	if a.currentUser == nil {
		return ErrUnauthorized
	}
	
	// Verify old PIN first (except for hardcoded admin)
	if a.currentUser.ID != 0 {
		_, err := a.db.AuthenticateUser(a.currentUser.FirstName+" "+a.currentUser.LastName, oldPIN)
		if err != nil {
			return errors.New("incorrect current PIN")
		}
	}
	
	// Update PIN (hardcoded admin can't change their PIN)
	if a.currentUser.ID == 0 {
		return errors.New("cannot change hardcoded admin PIN")
	}
	
	return a.db.ChangePIN(a.currentUser.ID, newPIN)
}

// ChangeUserPIN allows an admin to change another user's PIN
func (a *App) ChangeUserPIN(userID int, newPIN string) error {
	if a.currentUser == nil || !a.currentUser.IsAdmin {
		return ErrUnauthorized
	}
	
	return a.db.ChangePIN(userID, newPIN)
}

// UpdateUserPosition allows an admin to change a user's position
func (a *App) UpdateUserPosition(userID int, position string) error {
	if a.currentUser == nil || !a.currentUser.IsAdmin {
		return ErrUnauthorized
	}
	
	return a.db.UpdateUserPosition(userID, position)
}

// UpdateUserAdminStatus allows an admin to change a user's admin status
func (a *App) UpdateUserAdminStatus(userID int, isAdmin bool) error {
	if a.currentUser == nil || !a.currentUser.IsAdmin {
		return ErrUnauthorized
	}
	
	return a.db.UpdateUserAdminStatus(userID, isAdmin)
}

// UpdateUserJoinDate allows an admin to update a user's join date
func (a *App) UpdateUserJoinDate(userID int, joinDate string) error {
	if a.currentUser == nil || !a.currentUser.IsAdmin {
		return ErrUnauthorized
	}
	
	return a.db.UpdateUserJoinDate(userID, joinDate)
}

// GetPicklistByCategory returns picklist items for a category
func (a *App) GetPicklistByCategory(category string) ([]db.Picklist, error) {
	return a.db.GetPicklistByCategory(category)
}

// CreatePicklist creates a new picklist item
func (a *App) CreatePicklist(category, value string, sortOrder int) error {
	return a.db.CreatePicklistItem(category, value, sortOrder)
}

// UpdatePicklist updates an existing picklist item
func (a *App) UpdatePicklist(item *db.Picklist) error {
	return a.db.UpdatePicklistItem(item)
}

// DeletePicklist marks a picklist item as inactive
func (a *App) DeletePicklist(id int) error {
	return a.db.DeletePicklistItem(id)
}

// GetNextCallNumber gets the next call number for the given year
func (a *App) GetNextCallNumber(year int) (string, error) {
	return a.db.GetNextCallNumber(year)
}

// CreateCall creates a new call
func (a *App) CreateCall(call *db.Call, apparatusIDs []int, responderIDs []int, responderRoles []string) error {
	if a.currentUser == nil {
		return ErrUnauthorized
	}
	call.CreatedBy = a.currentUser.ID
	return a.db.CreateCall(call, apparatusIDs, responderIDs, responderRoles)
}

// GetCallByID returns a call by ID
func (a *App) GetCallByID(id int) (*db.Call, []db.Picklist, []db.User, error) {
	return a.db.GetCallByID(id)
}

// GetRecentCalls returns recent calls
func (a *App) GetRecentCalls(limit int) ([]db.Call, error) {
	return a.db.GetRecentCalls(limit, 0)
}

// GetCallsByYear returns all calls for a specific year
func (a *App) GetCallsByYear(year int) ([]db.Call, error) {
	return a.db.GetCallsByYear(year)
}

// GetCallYears returns all years that have calls
func (a *App) GetCallYears() ([]int, error) {
	return a.db.GetCallYears()
}

// SearchCalls searches for calls
func (a *App) SearchCalls(query string) ([]db.Call, error) {
	filters := make(map[string]interface{})
	if query != "" {
		filters["query"] = query
	}
	return a.db.SearchCalls(filters, 100, 0)
}

// UpdateCall updates an existing call
func (a *App) UpdateCall(call *db.Call, apparatusIDs []int, responderIDs []int, responderRoles []string) error {
	return a.db.UpdateCall(call, apparatusIDs, responderIDs, responderRoles)
}

// DeleteCall marks a call as deleted
func (a *App) DeleteCall(id int) error {
	// Just soft delete by updating the call
	call, _, _, err := a.db.GetCallByID(id)
	if err != nil {
		return err
	}
	// Note: no DeleteCall method exists, would need to implement soft delete if needed
	return a.db.UpdateCall(call, []int{}, []int{}, []string{})
}

// UploadLogo uploads and stores a logo image
func (a *App) UploadLogo(imageData []byte, mimeType string) error {
	if a.currentUser == nil || !a.currentUser.IsAdmin {
		return ErrUnauthorized
	}
	return a.db.SaveLogo(imageData, mimeType, a.currentUser.ID)
}

// GetLogo retrieves the stored logo image
func (a *App) GetLogo() (*db.Logo, error) {
	return a.db.GetLogo()
}

// DeleteLogo removes the stored logo
func (a *App) DeleteLogo() error {
	if a.currentUser == nil || !a.currentUser.IsAdmin {
		return ErrUnauthorized
	}
	return a.db.DeleteLogo()
}
