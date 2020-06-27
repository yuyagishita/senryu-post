package api

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/yu-yagishita/senryu-post/users"
)

// LoggingMiddleware はログを出力する
func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	Service
}

func (mw logmw) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.Service.Uppercase(s)
	return
}

func (mw logmw) Count(s string) (n int) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = mw.Service.Count(s)
	return
}

func (mw logmw) Login(username, password string) (users.User, error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Login",
			"username", username,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.Service.Login(username, password)
}

func (mw logmw) Register(username, email, password string) (string, error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Register",
			"username", username,
			"email", email,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.Service.Register(username, email, password)
}
