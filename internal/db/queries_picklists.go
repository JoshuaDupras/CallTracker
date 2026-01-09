package db

import (
	"database/sql"
)

// GetPicklistByCategory returns all active picklist items for a category
func (db *DB) GetPicklistByCategory(category string) ([]Picklist, error) {
	rows, err := db.Query(`
		SELECT id, category, value, sort_order, active
		FROM picklists 
		WHERE category = ? AND active = 1 
		ORDER BY sort_order, value
	`, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Picklist
	for rows.Next() {
		var item Picklist
		err := rows.Scan(&item.ID, &item.Category, &item.Value, &item.SortOrder, &item.Active)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// GetAllPicklistCategories returns all unique picklist categories
func (db *DB) GetAllPicklistCategories() ([]string, error) {
	rows, err := db.Query(`
		SELECT DISTINCT category 
		FROM picklists 
		ORDER BY category
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// GetPicklistsByCategoryForAdmin returns all picklist items for a category (including inactive for admin)
func (db *DB) GetPicklistsByCategoryForAdmin(category string) ([]Picklist, error) {
	rows, err := db.Query(`
		SELECT id, category, value, sort_order, active
		FROM picklists 
		WHERE category = ?
		ORDER BY sort_order, value
	`, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Picklist
	for rows.Next() {
		var item Picklist
		err := rows.Scan(&item.ID, &item.Category, &item.Value, &item.SortOrder, &item.Active)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// CreatePicklistItem creates a new picklist item
func (db *DB) CreatePicklistItem(category, value string, sortOrder int) error {
	_, err := db.Exec(`
		INSERT INTO picklists (category, value, sort_order, active) 
		VALUES (?, ?, ?, 1)
	`, category, value, sortOrder)
	return err
}

// UpdatePicklistItem updates a picklist item
func (db *DB) UpdatePicklistItem(item *Picklist) error {
	_, err := db.Exec(`
		UPDATE picklists 
		SET value = ?, sort_order = ?, active = ?
		WHERE id = ?
	`, item.Value, item.SortOrder, item.Active, item.ID)
	return err
}

// DeletePicklistItem soft-deletes a picklist item (sets active = false)
func (db *DB) DeletePicklistItem(id int) error {
	_, err := db.Exec(`
		UPDATE picklists 
		SET active = 0
		WHERE id = ?
	`, id)
	return err
}

// GetPicklistItem returns a single picklist item by ID
func (db *DB) GetPicklistItem(id int) (*Picklist, error) {
	var item Picklist
	err := db.QueryRow(`
		SELECT id, category, value, sort_order, active
		FROM picklists WHERE id = ?
	`, id).Scan(&item.ID, &item.Category, &item.Value, &item.SortOrder, &item.Active)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}
