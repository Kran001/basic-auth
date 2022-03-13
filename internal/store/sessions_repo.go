package store

import (
	"context"
	"fmt"

	"github.com/Kran001/basic-auth/internal/domain"

	"github.com/Kran001/basic-auth/pkg/database"

	"github.com/Kran001/basic-auth/pkg/logging"
	"github.com/jmoiron/sqlx"
)

type SessionsRepository struct {
	db *sqlx.DB
}

func NewSessionsRepository(db *sqlx.DB) Sessions {
	return &SessionsRepository{
		db: db,
	}
}

func (s *SessionsRepository) AddSession(ctx context.Context, session domain.Session) (string, error) {
	query := s.db.Rebind(`
		INSERT INTO 
		    sessions (token, user_id, user_agent, ip, session_time) 
		VALUES 
		    (:token, :user_id, :user_agent, :ip, :session_time)
		RETURNING id;
	`)

	token := session.Token
	err := database.WithTransaction(ctx, s.db, func(tx *sqlx.Tx) error {
		stmt, err := tx.PrepareNamedContext(ctx, query)
		if err != nil {
			return err
		}
		defer func() { _ = stmt.Close() }()

		_, err = stmt.ExecContext(ctx, session)

		return err
	})

	return token, err
}

func (s *SessionsRepository) AllUsersSessionsInfo(ctx context.Context) ([]domain.Session, error) {
	sessions := make([]domain.Session, 0)
	query := s.db.Rebind(`
		SELECT 
			s.id as session_id, s.token as token, u.id as user_id, 
		    u.name as name, u.psw as psw, 
		    u.mail as mail, u.params as params, s.user_agent as user_agent,
		    s.session_time as session_time, s.ip as ip
		FROM 
		     sessions s 
		INNER JOIN users u ON
		    u.id = s.user_id;
	`)

	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := conn.Close(); err != nil {
			logging.Logger.Error(errClose.Error())
		}
	}()

	if err = conn.SelectContext(ctx, &sessions, query); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (s *SessionsRepository) UserSessions(ctx context.Context, userId int64) (domain.Session, error) {
	session := domain.Session{}
	query := s.db.Rebind(`
		SELECT 
			s.id as session_id, s.token as token, u.id as user_id, 
		    u.name as name, u.psw as psw, 
		    u.mail as mail, u.params as params, s.user_agent as user_agent,
		    s.session_time as session_time, s.ip as ip
		FROM 
		     sessions s 
		INNER JOIN users u ON
		    u.id = s.user_id
		WHERE u.id = ?;
	`)

	conn, err := s.db.Connx(ctx)
	if err != nil {
		return domain.Session{}, err
	}

	defer func() {
		if errClose := conn.Close(); err != nil {
			logging.Logger.Error(errClose.Error())
		}
	}()

	if err = conn.GetContext(ctx, &session, query, userId); err != nil {
		return domain.Session{}, err
	}

	return session, nil
}

func (s *SessionsRepository) DeleteUserSessionById(ctx context.Context, sessionId int64) error {
	query := s.db.Rebind(`DELETE FROM sessions WHERE id = ?`)
	return database.WithTransaction(ctx, s.db, func(tx *sqlx.Tx) error {
		stmt, err := tx.PreparexContext(ctx, query)
		if err != nil {
			return err
		}
		defer func() { _ = stmt.Close() }()

		_, err = stmt.ExecContext(ctx, sessionId)

		return err
	})
}

func (s *SessionsRepository) DeleteAllUserSessions(ctx context.Context, userId int64) error {
	query := s.db.Rebind(`DELETE FROM sessions WHERE user_id = ?`)
	return database.WithTransaction(ctx, s.db, func(tx *sqlx.Tx) error {
		stmt, err := tx.PreparexContext(ctx, query)
		if err != nil {
			return nil
		}
		defer func() { _ = stmt.Close() }()

		_, err = stmt.ExecContext(ctx, userId)

		return err
	})
}

func (s *SessionsRepository) DeleteSessionByToken(ctx context.Context, token string) error {
	query := s.db.Rebind(`DELETE FROM sessions WHERE token = ?`)
	return database.WithTransaction(ctx, s.db, func(tx *sqlx.Tx) error {
		stmt, err := tx.PreparexContext(ctx, query)
		if err != nil {
			return nil
		}
		defer func() { _ = stmt.Close() }()

		_, err = stmt.ExecContext(ctx, token)

		return err
	})
}

func (s *SessionsRepository) CheckSession(ctx context.Context, spoilTimeMetric, spoilTime string) error {
	timeIntervalString := fmt.Sprintf("%s %s", spoilTime, spoilTimeMetric)
	query := s.db.Rebind(`
		DELETE FROM 
			sessions
		WHERE 
			date_trunc('?', now() - date_trunc( '?', session_time))	> '?'::interval;	
	`)

	return database.WithTransaction(ctx, s.db, func(tx *sqlx.Tx) error {
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		defer func() { _ = stmt.Close() }()

		_, err = stmt.ExecContext(ctx, spoilTimeMetric, spoilTimeMetric, timeIntervalString)

		return err
	})
}
