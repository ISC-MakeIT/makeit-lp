package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

const mysqlDuplicateEntry = 1062

func findOrCreateGoogleUser(
	ctx context.Context,
	db *sql.DB,
	identity GoogleIdentity,
) (int64, error) {
	if len(identity.Subject) == 0 || len(identity.Subject) > 255 {
		return 0, errors.New("invalid verified Google subject length")
	}

	userID, err := findUserIDByGoogleSubject(ctx, db, identity.Subject)
	if err == nil {
		return userID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("look up user by Google subject: %w", err)
	}

	const insertUser = `INSERT INTO users (google_sub) VALUES (?)`
	insertedUser, insertErr := db.ExecContext(
		ctx, insertUser, []byte(identity.Subject),
	)
	if insertErr == nil {
		userID, err := insertedUser.LastInsertId()
		if err != nil {
			return 0, fmt.Errorf("read new site user ID: %w", err)
		}
		return userID, nil
	}

	var mysqlErr *mysql.MySQLError
	if !errors.As(insertErr, &mysqlErr) || mysqlErr.Number != mysqlDuplicateEntry {
		return 0, fmt.Errorf("insert site user: %w", insertErr)
	}

	// 別の列の重複でも1062になるため、同じGoogle識別子で検索し直す。
	userID, err = findUserIDByGoogleSubject(ctx, db, identity.Subject)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("insert site user: %w", insertErr)
	}
	if err != nil {
		return 0, fmt.Errorf("look up concurrent Google registration: %w", err)
	}

	return userID, nil
}

func findUserIDByGoogleSubject(
	ctx context.Context,
	db *sql.DB,
	subject string,
) (int64, error) {
	const query = `SELECT id FROM users WHERE google_sub = ?`

	var userID int64
	if err := db.QueryRowContext(ctx, query, []byte(subject)).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}
