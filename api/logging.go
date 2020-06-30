package api

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/yu-yagishita/senryu-post/posts"
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

func (mw logmw) GetAll() (posts.Post, error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetAll",
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.Service.GetAll()
}
