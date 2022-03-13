package dispatchers

import (
	"context"
	"time"

	"github.com/Kran001/basic-auth/internal/store"
	"github.com/Kran001/basic-auth/pkg/logging"
)

const CheckDelay = 10 * time.Minute

type SessionCleaner struct {
	db              store.Sessions
	checkTicker     *time.Ticker
	spoilTime       string
	spoilTimeMetric string
}

func NewSessionCleaner(db store.Sessions, spoilTime, spoilTimeMetric string) *SessionCleaner {
	return &SessionCleaner{
		db:              db,
		checkTicker:     time.NewTicker(CheckDelay),
		spoilTime:       spoilTime,
		spoilTimeMetric: spoilTimeMetric,
	}
}

func (s *SessionCleaner) Run(ctx context.Context) error {
	defer func() {
		s.checkTicker.Stop()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.checkTicker.C:
		logging.Logger.Info("Checking sessions targets for cleanup")

		if err := s.db.CheckSession(ctx, s.spoilTime, s.spoilTimeMetric); err != nil {
			logging.Logger.Error("Failed check session options: ", err.Error())
		}
	}

	return nil
}
