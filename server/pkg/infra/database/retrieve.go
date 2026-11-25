package database

import (
	"database/sql"
	"errors"
)

// MustExec checks the result of a SQL execution.
// It returns an error if no rows were affected.
func MustExec(result sql.Result, err error) (sql.Result, error) {
	if err != nil {
		return result, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return result, err
	}
	// Check if at least one row was affected
	if rowsAffected == 0 {
		return result, errors.New("rowsAffected is 0")
	}

	return result, nil
}
