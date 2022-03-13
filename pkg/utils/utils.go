package utils

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Kran001/basic-auth/pkg/logging"
)

// SimpleError Common error struct (for sending to web).
type SimpleError struct {
	Msg  string // message text
	Code int    // error code
}

// NotifyChan Linking in chan to os terminate signals.
func NotifyChan(in chan os.Signal) {
	signal.Notify(in, os.Interrupt, syscall.SIGTERM)
	signal.Notify(in, os.Interrupt, syscall.SIGINT)
}

func WaitSignals(ctx context.Context) error {
	sig := make(chan os.Signal, 1)
	NotifyChan(sig)

	select {
	case <-ctx.Done():
	case s := <-sig:
		logging.Logger.Info(s)
	}

	return nil
}

func MakeErr(err string) SimpleError {
	var nErr SimpleError
	twoStr := strings.Split(err, ", message: ")
	coStr := strings.Split(twoStr[0], "code: ")

	logging.Logger.Info(twoStr, coStr)

	if len(coStr) > 1 {
		nErr.Code, _ = strconv.Atoi(coStr[1])
	} else {
		nErr.Code = 0
	}

	if len(twoStr) > 1 {
		nErr.Msg = twoStr[1]
	} else {
		if err != "" {
			nErr.Msg = err
		} else {
			nErr.Msg = "Smth is wrong..."
		}
	}

	return nErr
}
