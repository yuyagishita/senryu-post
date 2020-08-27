package api

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/yu-yagishita/senryu-post/posts"
)

// InstrumentingMiddleware はアクセスの計測ができる
func InstrumentingMiddleware(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
	countResult metrics.Histogram,
) ServiceMiddleware {
	return func(next Service) Service {
		return instrmw{requestCount, requestLatency, countResult, next}
	}
}

type instrmw struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	Service
}

func (mw instrmw) GetAll() (post []posts.Post, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetAll", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Service.GetAll()
}
