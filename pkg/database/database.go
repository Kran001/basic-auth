package database

import (
	"context"
	"database/sql"

	"github.com/Kran001/basic-auth/pkg/logging"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/lib/pq"
)

type Config interface {
	ConnectionString() string
}

func NewDBConnection(c Config) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(context.Background(), "postgres", c.ConnectionString())
	if err != nil {
		logging.Logger.Error("Error connector to database with error:", err.Error())

		return nil, err
	}

	db.Mapper = reflectx.NewMapperFunc("json", func(s string) string { return s })

	return db, nil
}

type TxFn func(*sqlx.Tx) error

func WithTransaction(ctx context.Context, db *sqlx.DB, fx TxFn) error {
	conn, err := db.Connx(ctx)
	if err != nil {
		logging.Logger.Error("Error get sqlx connection. Reason:", err.Error())

		return err
	}

	defer func() {
		if errClose := conn.Close(); errClose != nil {
			logging.Logger.Error("Failed close conn. Reason: ", errClose.Error())
		}
	}()

	tx, err := conn.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		logging.Logger.Error("Error start transaction: ", err.Error())

		return err
	}

	defer func() {
		if p := recover(); p != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				logging.Logger.Error(errRollback.Error())
			}
		} else if err != nil {
			// something went wrong, rollback
			if errRollback := tx.Rollback(); errRollback != nil {
				logging.Logger.Error("Something wrong:", errRollback.Error())
			}
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	return fx(tx)
}
