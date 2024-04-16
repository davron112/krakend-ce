package krakend

import (
	"context"

	"api-gateway/v2/modules/lura/v2/config"
	"api-gateway/v2/modules/lura/v2/logging"
	"api-gateway/v2/modules/lura/v2/sd/dnssrv"
)

// RegisterSubscriberFactories registers all the available sd adaptors
func RegisterSubscriberFactories(_ context.Context, _ config.ServiceConfig, _ logging.Logger) func(n string, p int) {
	// register the dns service discovery
	dnssrv.Register()

	return func(name string, port int) {}
}

type registerSubscriberFactories struct{}

func (registerSubscriberFactories) Register(ctx context.Context, cfg config.ServiceConfig, logger logging.Logger) func(n string, p int) {
	return RegisterSubscriberFactories(ctx, cfg, logger)
}
