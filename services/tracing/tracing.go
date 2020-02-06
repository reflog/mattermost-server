// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package tracing

import (
	"io"
	"time"

	"github.com/mattermost/mattermost-server/v5/mlog"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"

	"context"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type Tracer struct {
	closer io.Closer
}

type LogrusAdapter struct {
}

// Error - logrus adapter for span errors
func (l LogrusAdapter) Error(msg string) {
	mlog.Error(msg)
}

// Infof - logrus adapter for span info logging
func (l LogrusAdapter) Infof(msg string, args ...interface{}) {
	// we ignore Info messages from opentracing
}

func New() (*Tracer, error) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	closer, err := cfg.InitGlobalTracer(
		"mattermost",
		jaegercfg.Logger(LogrusAdapter{}),
		jaegercfg.Metrics(metrics.NullFactory),
		jaegercfg.Tag("serverStartTime", time.Now().UTC().Format(time.RFC3339)),
	)

	if err != nil {
		return nil, err
	}
	mlog.Info("Opentracing initialzed")
	return &Tracer{
		closer: closer,
	}, nil
}

func (t *Tracer) Close() error {
	return t.closer.Close()
}

func StartRootSpanByContext(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, operationName)
}

func StartSpanWithParentByContext(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	parentSpan := opentracing.SpanFromContext(ctx)

	if parentSpan == nil {
		return StartRootSpanByContext(ctx, operationName)
	}

	return opentracing.StartSpanFromContext(ctx, operationName, opentracing.ChildOf(parentSpan.Context()))
}