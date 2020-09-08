package mimir

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

func Tracer(service, version string, logger Logging) (opentracing.Tracer, func(), error) {
	cfg, err := jaegerCfg.FromEnv()
	if err != nil {
		logger.Errorf("Could not parse Jaeger env vars: %s", err.Error())
		return nil, nil, err
	}

	cfg.Sampler.Param = 1

	if service != "" {
		cfg.ServiceName = service
	}
	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.NewTracer(
		jaegerCfg.Metrics(jMetricsFactory),
		jaegerCfg.Logger(logger),
		jaegerCfg.Tag(fmt.Sprintf("%.version", service), version),
		jaegerCfg.MaxTagValueLength(2048),
	)
	if err != nil {
		logger.Errorf("Could not initialize jaeger tracer: %s", err.Error())
		return nil, nil, err
	}

	cleanup := func() {
		_ = closer.Close()
	}

	return tracer, cleanup, err
}

func TracerSpanCallback(ctx context.Context, operation string, f func(span context.Context) error) error {
	span, ctxSpan := opentracing.StartSpanFromContext(ctx, operation)
	defer span.Finish()
	return f(ctxSpan)
}
