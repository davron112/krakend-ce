package krakend

import (
	"context"
	"fmt"
	amqp "api-gateway/v2/modules/krakend-amqp/v2"
	cel "api-gateway/v2/modules/krakend-cel/v2"
	cb "api-gateway/v2/modules/krakend-circuitbreaker/v2/gobreaker/proxy"
	httpcache "api-gateway/v2/modules/krakend-httpcache/v2"
	lambda "api-gateway/v2/modules/krakend-lambda/v2"
	lua "api-gateway/v2/modules/krakend-lua/v2/proxy"
	martian "api-gateway/v2/modules/krakend-martian/v2"
	metrics "api-gateway/v2/modules/krakend-metrics/v2/gin"
	oauth2client "api-gateway/v2/modules/krakend-oauth2-clientcredentials/v2"
	opencensus "api-gateway/v2/modules/krakend-opencensus/v2"
	pubsub "api-gateway/v2/modules/krakend-pubsub/v2"
	ratelimit "api-gateway/v2/modules/krakend-ratelimit/v3/proxy"
	"api-gateway/v2/modules/lura/v2/config"
	"api-gateway/v2/modules/lura/v2/logging"
	"api-gateway/v2/modules/lura/v2/proxy"
	"api-gateway/v2/modules/lura/v2/transport/http/client"
	httprequestexecutor "api-gateway/v2/modules/lura/v2/transport/http/client/plugin"
	"api-gateway/v2/proxy/grpc"
)

// NewBackendFactory creates a BackendFactory by stacking all the available middlewares:
// - oauth2 client credentials
// - http cache
// - martian
// - pubsub
// - amqp
// - cel
// - lua
// - rate-limit
// - circuit breaker
// - metrics collector
// - opencensus collector
func NewBackendFactory(logger logging.Logger, metricCollector *metrics.Metrics) proxy.BackendFactory {
	return NewBackendFactoryWithContext(context.Background(), logger, metricCollector)
}

// NewBackendFactoryWithContext creates a BackendFactory by stacking all the available middlewares and injecting the received context
func NewBackendFactoryWithContext(ctx context.Context, logger logging.Logger, metricCollector *metrics.Metrics) proxy.BackendFactory {
	requestExecutorFactory := func(cfg *config.Backend) client.HTTPRequestExecutor {
		clientFactory := client.NewHTTPClient
		if _, ok := cfg.ExtraConfig[oauth2client.Namespace]; ok {
			clientFactory = oauth2client.NewHTTPClient(cfg)
		} else {
			clientFactory = httpcache.NewHTTPClient(cfg, clientFactory)
		}
		return opencensus.HTTPRequestExecutorFromConfig(clientFactory, cfg)
	}
	requestExecutorFactory = httprequestexecutor.HTTPRequestExecutor(logger, requestExecutorFactory)
	backendFactory := martian.NewConfiguredBackendFactory(logger, requestExecutorFactory)
	bf := pubsub.NewBackendFactory(ctx, logger, backendFactory)
	backendFactory = bf.New
	backendFactory = grpc.NewGrpcBackendFactory(logger, backendFactory)
	backendFactory = amqp.NewBackendFactory(ctx, logger, backendFactory)
	backendFactory = lambda.BackendFactory(logger, backendFactory)
	backendFactory = cel.BackendFactory(logger, backendFactory)
	backendFactory = lua.BackendFactory(logger, backendFactory)
	backendFactory = ratelimit.BackendFactory(logger, backendFactory)
	backendFactory = cb.BackendFactory(backendFactory, logger)
	backendFactory = metricCollector.BackendFactory("backend", backendFactory)
	backendFactory = opencensus.BackendFactory(backendFactory)

	return func(remote *config.Backend) proxy.Proxy {
		logger.Debug(fmt.Sprintf("[BACKEND: %s] Building the backend pipe", remote.URLPattern))
		return backendFactory(remote)
	}
}

type backendFactory struct{}

func (backendFactory) NewBackendFactory(ctx context.Context, l logging.Logger, m *metrics.Metrics) proxy.BackendFactory {
	return NewBackendFactoryWithContext(ctx, l, m)
}
