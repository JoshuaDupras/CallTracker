package db

import (
	"time"
)

// User represents a fire department member
type User struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Position   string    `json:"position"` // "chief", "captain", "member", "probationary"
	EMSLevel   string    `json:"ems_level"` // "VEFR", "EMR", "EMT", "AEMT", "Paramedic", "None"
	IsAdmin    bool      `json:"is_admin"`
	PIN        string    `json:"pin,omitempty"`
	Active     bool      `json:"active"`
	JoinedDate *time.Time `json:"joined_date,omitempty"`
	Created    time.Time `json:"created"`
}

// Picklist represents dropdown values for various categories
type Picklist struct {
	ID        int    `json:"id"`
	Category  string `json:"category"` // town, call_type, priority, etc.
	Value     string `json:"value"`
	SortOrder int    `json:"sort_order"`
	Active    bool   `json:"active"`
}

// FormField represents form configuration
type FormField struct {
	ID        int    `json:"id"`
	FieldName string `json:"field_name"`
	Label     string `json:"label"`
	Required  bool   `json:"required"`
	Enabled   bool   `json:"enabled"`
	SortOrder int    `json:"sort_order"`
}

// Call represents a fire department call
type Call struct {
	ID             int       `json:"id"`
	IncidentNumber string    `json:"incident_number"`
	CallType       string    `json:"call_type"`
	MutualAid      string    `json:"mutual_aid"`
	Address        string    `json:"address"`
	Town           string    `json:"town"`
	LocationNotes  string    `json:"location_notes"`
	Dispatched     time.Time `json:"dispatched"`
	Enroute        *time.Time `json:"enroute"`
	OnScene        *time.Time `json:"on_scene"`
	Clear          *time.Time `json:"clear"`
	Narrative      string    `json:"narrative"`
	CreatedBy      int       `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CallApparatus represents apparatus assigned to a call
type CallApparatus struct {
	ID          int `json:"id"`
	CallID      int `json:"call_id"`
	ApparatusID int `json:"apparatus_id"`
}

// CallResponder represents responders assigned to a call
type CallResponder struct {
	ID            int    `json:"id"`
	CallID        int    `json:"call_id"`
	ResponderID   int    `json:"responder_id"`
	ResponderRole string `json:"responder_role,omitempty"`
}

// Setting represents application configuration
type Setting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// AuditLog represents audit trail for accountability
type AuditLog struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Action    string    `json:"action"`
	TableName string    `json:"table_name"`
	RecordID  int       `json:"record_id"`
	Changes   string    `json:"changes"`
	Timestamp time.Time `json:"timestamp"`
}

// Logo represents the uploaded logo image
type Logo struct {
	ID         int       `json:"id"`
	ImageData  []byte    `json:"image_data"`
	MimeType   string    `json:"mime_type"`
	UploadedAt time.Time `json:"uploaded_at"`
	UploadedBy int       `json:"uploaded_by"`
}
