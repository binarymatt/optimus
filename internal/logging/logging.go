package logging

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel/trace"
)

// Unexported new type so that our context key never collides with another.
type contextKeyType struct{}

// contextKey is the key used for the context to store the logger.
var contextKey = contextKeyType{}

func WithContext(ctx context.Context, logger *slog.Logger, args ...interface{}) context.Context {
	// While we could call logger.With even with zero args, we have this
	// check to avoid unnecessary allocations around creating a copy of a
	// logger.
	if len(args) > 0 {
		logger = logger.With(args...)
	}

	return context.WithValue(ctx, contextKey, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	logger, _ := ctx.Value(contextKey).(*slog.Logger)
	if logger == nil {
		return slog.Default()
	}

	return logger
}

func NewInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			httpMethod := req.HTTPMethod()
			method := req.Spec().Procedure
			scontext := trace.SpanContextFromContext(ctx)
			logger := slog.Default()
			if scontext.HasTraceID() {
				logger = logger.With("trace_id", scontext.TraceID())
			}
			WithContext(ctx, logger)
			logger.DebugContext(ctx, "grpc endpoint called", slog.Group("grpc", slog.String("http_method", httpMethod), slog.String("method", method)))
			//slog.SetDefault(slog.Default().With(slog.Group("grpc", slog.String("http_method", httpMethod), slog.String("method", method))))
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
