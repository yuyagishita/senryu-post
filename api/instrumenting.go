package api

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/yu-yagishita/senryu-post/users"
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

func (mw instrmw) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "uppercase", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.Service.Uppercase(s)
	return
}

func (mw instrmw) Count(s string) (n int) {
	defer func(begin time.Time) {
		lvs := []string{"method", "count", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mw.countResult.Observe(float64(n))
	}(time.Now())

	n = mw.Service.Count(s)
	return
}

func (mw instrmw) Login(username, password string) (user users.User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Login", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Service.Login(username, password)
}

func (mw instrmw) Register(username, email, password string) (str string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Register", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Service.Register(username, email, password)
}
