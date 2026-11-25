package database

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ExecExecutor is the minimal interface required for non-transactional updates.
type ExecExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// BuildUpdateQuery 根据传入结构体中非 nil 的指针字段，动态生成 UPDATE SQL 语句及其对应的参数。
// 返回生成的 SQL 语句和参数切片。
func BuildUpdateQuery(tableName string, updates any, whereClause string, whereArgs ...any) (string, []any, error) {
	v := reflect.ValueOf(updates)
	t := reflect.TypeOf(updates)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("updates must be a struct or pointer to struct")
	}

	var setClauses []string
	var args []any
	argIndex := 1

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typeField := t.Field(i)

		dbTag := typeField.Tag.Get("db")
		if dbTag == "" || dbTag == "-" {
			continue
		}

		if field.Kind() != reflect.Ptr {
			continue
		}

		if field.IsNil() {
			continue
		}

		value := field.Elem().Interface()

		setClauses = append(setClauses, fmt.Sprintf("%s=$%d", dbTag, argIndex))
		args = append(args, value)
		argIndex++
	}

	if len(setClauses) == 0 {
		return "", nil, fmt.Errorf("no fields to update")
	}

	// Adjust placeholders in the WHERE clause so they continue after the SET args.
	// We only modify the WHERE clause portion to avoid accidentally changing
	// placeholders in the SET clause.
	if len(whereArgs) > 0 {
		// Use regex to replace placeholders safely
		// Matches $ followed by digits
		re := regexp.MustCompile(`\$(\d+)`)

		adjustedWhere := re.ReplaceAllStringFunc(whereClause, func(match string) string {
			// Extract the number part
			numStr := match[1:]
			idx, err := strconv.Atoi(numStr)
			if err != nil {
				return match // Should not happen given the regex
			}

			// Add the offset (len(args))
			newIdx := idx + len(args)
			return fmt.Sprintf("$%d", newIdx)
		})

		query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, strings.Join(setClauses, ", "), adjustedWhere)
		finalArgs := make([]any, 0, len(args)+len(whereArgs))
		finalArgs = append(finalArgs, args...)
		finalArgs = append(finalArgs, whereArgs...)
		return query, finalArgs, nil
	}

	// No where args: simple query
	query := fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(setClauses, ", "))
	return query, args, nil
}

// ExecuteUpdate executes an update query using BuildUpdateQuery with the given db connection.
func ExecuteUpdate(db ExecExecutor, tableName string, updates any, whereClause string, whereArgs ...any) error {
	query, args, err := BuildUpdateQuery(tableName, updates, whereClause, whereArgs...)
	if err != nil {
		return err
	}

	_, err = db.ExecContext(context.Background(), query, args...)
	return err
}

// ExecuteUpdateTx executes an update query within a transaction using BuildUpdateQuery
func ExecuteUpdateTx(tx *sql.Tx, tableName string, updates any, whereClause string, whereArgs ...any) error {
	query, args, err := BuildUpdateQuery(tableName, updates, whereClause, whereArgs...)
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, args...)
	return err
}
