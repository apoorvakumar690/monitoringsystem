// Package apmgrpc provides interceptors for tracing monitoring gRPC.
package apmgrpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/apoorvakumar690/monitoringsystem"
)

var (
	defaultOptions = &options{}
)

// key Context apm transaction key
type key string

const (
	// ElasticConextKey Use context key to get values only for request-scoped data that transits processes and
	// APIs, not for passing optional parameters to functions.
	ElasticConextKey key = "apm"
)

// options options for creating a request context object
type options struct {
	apm *monitoringsystem.Agent
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Option sets options for server-side tracing.
type Option func(*options)

// WithAPM customizes the function for monitoring the request performance
func WithAPM(apm *monitoringsystem.Agent) Option {
	return func(o *options) {
		o.apm = apm
	}
}

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor that
// traces gRPC requests with the given options.
//
// The interceptor will trace transactions with the "grpc" type for each
// incoming request.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if o.apm != nil {
			tx, err := o.apm.StartTransaction(info.FullMethod)
			o.apm.AddAttribute(tx, "gRPC", grpc.Version)
			ctx = context.WithValue(ctx, ElasticConextKey, tx)
			defer o.apm.EndTransaction(tx, err)
		}
		resp, err = handler(ctx, req)
		return resp, err
	}
}

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor that
// traces gRPC requests with the given apm options.
func UnaryClientInterceptor(opts ...Option) grpc.UnaryClientInterceptor {
	o := evaluateOptions(opts)
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if o.apm != nil {
			tx, err := o.apm.StartTransaction(method)
			o.apm.AddAttribute(tx, "gRPC", grpc.Version)
			ctx = context.WithValue(ctx, ElasticConextKey, tx)
			defer o.apm.EndTransaction(tx, err)
		}
		return invoker(ctx, method, req, resp, cc, opts...)
	}
}
