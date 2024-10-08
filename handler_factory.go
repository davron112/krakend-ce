package krakend

import (
	botdetector "api-gateway/v2/modules/krakend-botdetector/v2/gin"
	jose "api-gateway/v2/modules/krakend-jose/v2"
	ginjose "api-gateway/v2/modules/krakend-jose/v2/gin"
	lua "api-gateway/v2/modules/krakend-lua/v2/router/gin"
	metrics "api-gateway/v2/modules/krakend-metrics/v2/gin"
	opencensus "api-gateway/v2/modules/krakend-opencensus/v2/router/gin"
	ratelimit "api-gateway/v2/modules/krakend-ratelimit/v3/router/gin"
	"api-gateway/v2/modules/lura/v2/config"
	"api-gateway/v2/modules/lura/v2/logging"
	"api-gateway/v2/modules/lura/v2/proxy"
	router "api-gateway/v2/modules/lura/v2/router/gin"
	"api-gateway/v2/modules/lura/v2/transport/http/server"
	"api-gateway/v2/pkg/jwtvalidator"
	"fmt"

	"github.com/gin-gonic/gin"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(logger logging.Logger, metricCollector *metrics.Metrics, rejecter jose.RejecterFactory) router.HandlerFactory {
	handlerFactory := router.CustomErrorEndpointHandler(logger, server.DefaultToHTTPError)
	handlerFactory = ratelimit.NewRateLimiterMw(logger, handlerFactory)
	handlerFactory = lua.HandlerFactory(logger, handlerFactory)
	handlerFactory = ginjose.HandlerFactory(handlerFactory, logger, rejecter)
	handlerFactory = metricCollector.NewHTTPHandlerFactory(handlerFactory)
	handlerFactory = opencensus.New(handlerFactory)
	handlerFactory = botdetector.New(handlerFactory, logger)
	handlerFactory = jwtvalidator.New(handlerFactory, logger)

	return func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		logger.Debug(fmt.Sprintf("[ENDPOINT: %s] Building the http handler", cfg.Endpoint))
		return handlerFactory(cfg, p)
	}
}

type handlerFactory struct{}

func (handlerFactory) NewHandlerFactory(l logging.Logger, m *metrics.Metrics, r jose.RejecterFactory) router.HandlerFactory {
	return NewHandlerFactory(l, m, r)
}
