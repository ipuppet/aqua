// Package types holds shared value types used across domain models.
package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullString is a JSON-friendly sql.NullString for optional string columns.
type NullString struct {
	sql.NullString
}

// Value implements sql/driver.Valuer. Invalid values are stored as NULL.
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// Scan implements sql.Scanner.
func (ns *NullString) Scan(src any) error {
	return ns.NullString.Scan(src)
}

// MarshalJSON encodes invalid values as JSON null.
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON accepts JSON null as an invalid value.
func (ns *NullString) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ns.String, ns.Valid = "", false
		return nil
	}
	if err := json.Unmarshal(b, &ns.String); err != nil {
		return err
	}
	ns.Valid = true
	return nil
}
