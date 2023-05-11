package ddtracer

import (
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const componentName = "labstack/echo.v4"

func Middleware(opts ...Option) echo.MiddlewareFunc {
	cfg := new(config)
	defaults(cfg)
	for _, fn := range opts {
		fn(cfg)
	}
	spanOpts := []ddtrace.StartSpanOption{
		tracer.ServiceName(cfg.serviceName),
		tracer.Tag(ext.Component, componentName),
		tracer.Tag(ext.SpanKind, ext.SpanKindServer),
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			route := c.Path()
			resource := request.Method + " " + route
			opts := append(spanOpts, tracer.ResourceName(resource), tracer.Tag(ext.HTTPRoute, route))

			if !math.IsNaN(cfg.analyticsRate) {
				opts = append(opts, tracer.Tag(ext.EventSampleRate, cfg.analyticsRate))
			}

			var finishOpts []tracer.FinishOption
			if cfg.noDebugStack {
				finishOpts = []tracer.FinishOption{tracer.NoDebugStack()}
			}

			span, ctx := StartRequestSpan(request, opts...)
			defer func() {
				FinishRequestSpan(span, c.Response().Status, finishOpts...)
			}()

			c.SetRequest(request.WithContext(ctx))

			err := next(c)

			if err != nil {
				c.Error(err)
				finishOpts = append(finishOpts, tracer.WithError(err))
			}

			span.SetTag(ext.HTTPCode, strconv.Itoa(c.Response().Status))

			return err
		}
	}
}
