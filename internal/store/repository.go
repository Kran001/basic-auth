package store

import (
	"context"

	"github.com/Kran001/basic-auth/internal/domain"

	"github.com/jmoiron/sqlx"
)

type Users interface {
	AddNewUser(ctx context.Context, user domain.User) (int64, error)
	UserById(ctx context.Context, id int64) (domain.User, error)
	UserByName(ctx context.Context, name string) (domain.User, error)
	UsersList(ctx context.Context) ([]domain.User, error)
	FindUser(ctx context.Context, email string, password string) (domain.User, error)
}

type Sessions interface {
	AddSession(ctx context.Context, session domain.Session) (string, error)
	AllUsersSessionsInfo(ctx context.Context) ([]domain.Session, error)
	UserSessions(ctx context.Context, userId int64) (domain.Session, error)
	DeleteUserSessionById(ctx context.Context, sessionId int64) error
	DeleteAllUserSessions(ctx context.Context, userId int64) error
	DeleteSessionByToken(ctx context.Context, token string) error
	CheckSession(ctx context.Context, spoilTimeMetric, spoilTime string) error
}

type Repositories struct {
	Users    Users
	Sessions Sessions
}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		Users:    NewUsersRepository(db),
		Sessions: NewSessionsRepository(db),
	}
}
