package store

import (
	"context"
	"errors"
	"strings"

	"github.com/Kran001/basic-auth/internal/domain"

	"github.com/Kran001/basic-auth/pkg/database"

	"github.com/Kran001/basic-auth/pkg/logging"
	"github.com/jmoiron/sqlx"
)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) Users {
	return &UsersRepository{
		db: db,
	}
}

func (u *UsersRepository) AddNewUser(ctx context.Context, user domain.User) (int64, error) {
	id := int64(0)
	query := u.db.Rebind(`
		INSERT INTO 
		    users (name, password, email) 
		VALUES 
		    (:name, :password, :email) 
		RETURNING id;
	`)

	err := database.WithTransaction(ctx, u.db, func(tx *sqlx.Tx) error {
		stmt, err := tx.PrepareNamedContext(ctx, query)
		if err != nil {
			return err
		}
		defer func() { _ = stmt.Close() }()

		err = stmt.QueryRowxContext(ctx, user).Scan(&id)
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key") {
				return err
			}

			return errors.New("duplicate value")
		}

		return nil
	})

	return id, err
}

func (u *UsersRepository) UserById(ctx context.Context, id int64) (domain.User, error) {
	user := domain.User{}
	query := u.db.Rebind(`
		SELECT 
			u.id as user_id, u.name as name, u.password as password, u.mail as mail
		FROM users u
		WHERE u.id = ?;
	`)

	conn, err := u.db.Connx(ctx)
	if err != nil {
		return domain.User{}, err
	}

	defer func() {
		if errClose := conn.Close(); err != nil {
			logging.Logger.Error(errClose.Error())
		}
	}()

	if err = conn.GetContext(ctx, &user, query, id); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (u *UsersRepository) UserByName(ctx context.Context, name string) (domain.User, error) {
	user := domain.User{}
	query := u.db.Rebind(`
		SELECT 
			u.id as user_id, u.name as name, u.password as password, u.email as email
		FROM users u
		WHERE u.name = ?;
	`)

	conn, err := u.db.Connx(ctx)
	if err != nil {
		return domain.User{}, err
	}

	defer func() {
		if errClose := conn.Close(); err != nil {
			logging.Logger.Error(errClose.Error())
		}
	}()

	if err = conn.GetContext(ctx, &user, query, name); err != nil {
		return domain.User{}, err
	}

	return user, nil

}

func (u *UsersRepository) UsersList(ctx context.Context) ([]domain.User, error) {
	usersList := make([]domain.User, 0)
	query := u.db.Rebind(`
		SELECT 
			u.id as user_id, u.name as name, u.password as password, u.email as email
		FROM users u;
	`)

	conn, err := u.db.Connx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := conn.Close(); err != nil {
			logging.Logger.Error(errClose.Error())
		}
	}()

	if err = conn.SelectContext(ctx, &usersList, query); err != nil {
		return nil, err
	}

	return usersList, nil
}

func (u *UsersRepository) FindUser(ctx context.Context, email string, password string) (domain.User, error) {
	user := domain.User{}
	query := u.db.Rebind(`
		SELECT 
			u.id as user_id, u.name as name, u.password as password, u.email as email
		FROM users u
		WHERE u.email = ? AND u.password = ?;
	`)

	conn, err := u.db.Connx(ctx)
	if err != nil {
		return domain.User{}, err
	}

	defer func() {
		if errClose := conn.Close(); err != nil {
			logging.Logger.Error(errClose.Error())
		}
	}()

	if err = conn.GetContext(ctx, &user, query, email, password); err != nil {
		return domain.User{}, err
	}

	return user, nil
}
