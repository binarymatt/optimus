package service

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/protobuf/proto"
)

var (
	DefaultBytesDistribution        = []float64{1024, 2048, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864, 268435456, 1073741824, 4294967296}
	DefaultMillisecondsDistribution = []float64{0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200, 250, 300, 400, 500, 650, 800, 1000, 2000, 5000, 10000, 20000, 50000, 100000}
	DefaultMicrosecondsDistribution = []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000}
	DefaultMessageCountDistribution = []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536}

	handlerLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "handler_latency",
		Help:    "endpoint latency measured in milliseconds",
		Buckets: DefaultMillisecondsDistribution,
	}, []string{
		"endpoint",
		"status_code",
	})
	handlerRequestSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_request_size",
		Help:    "request size bytes",
		Buckets: DefaultBytesDistribution,
	}, []string{
		"endpoint",
	})
	handlerResponseSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "handler_response_size",
		Help:    "response size bytes",
		Buckets: DefaultBytesDistribution,
	}, []string{
		"endpoint",
		"status_code",
	})

	handlerCompleted = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_handled",
	}, []string{
		"endpoint",
		"status_code",
	})
)

func MetricsInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			start := time.Now()
			endpoint := req.Spec().Procedure

			var reqSize int
			if req != nil {
				if msg, ok := req.Any().(proto.Message); ok {
					reqSize = proto.Size(msg)
				}
			}
			resp, err := next(ctx, req)
			var respSize int
			statusCode := "OK"
			if err == nil {
				if msg, ok := resp.Any().(proto.Message); ok {
					respSize = proto.Size(msg)
				}
			} else {
				statusCode = connect.CodeOf(err).String()
			}
			duration := time.Since(start).Milliseconds()

			handlerRequestSize.WithLabelValues(endpoint).Observe(float64(reqSize))
			handlerResponseSize.WithLabelValues(endpoint, statusCode).Observe(float64(respSize))
			handlerLatency.WithLabelValues(endpoint, statusCode).Observe(float64(duration))
			handlerCompleted.WithLabelValues(endpoint, statusCode).Inc()

			return resp, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
