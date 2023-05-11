package ddtracer

import (
	"math"

	"github.com/labstack/echo/v4"
)

const defaultServiceName = "echo"

type config struct {
	ignoreRequestFunc IgnoreRequestFunc
	isStatusError     func(statusCode int) bool
	serviceName       string
	analyticsRate     float64
	noDebugStack      bool
}

type Option func(*config)

type IgnoreRequestFunc func(c echo.Context) bool

func defaults(cfg *config) {
	cfg.serviceName = "echo"
	cfg.analyticsRate = math.NaN()
	cfg.isStatusError = isServerError
}

func WithServiceName(name string) Option {
	return func(cfg *config) {
		cfg.serviceName = name
	}
}

func WithAnalytics(on bool) Option {
	return func(cfg *config) {
		if on {
			cfg.analyticsRate = 1.0
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}

func WithAnalyticsRate(rate float64) Option {
	return func(cfg *config) {
		if rate >= 0.0 && rate <= 1.0 {
			cfg.analyticsRate = rate
		} else {
			cfg.analyticsRate = math.NaN()
		}
	}
}

func NoDebugStack() Option {
	return func(cfg *config) {
		cfg.noDebugStack = true
	}
}

func WithIgnoreRequest(ignoreRequestFunc IgnoreRequestFunc) Option {
	return func(cfg *config) {
		cfg.ignoreRequestFunc = ignoreRequestFunc
	}
}

func WithStatusCheck(fn func(statusCode int) bool) Option {
	return func(cfg *config) {
		cfg.isStatusError = fn
	}
}

func isServerError(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}
