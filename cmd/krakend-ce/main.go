package main

import (
	"context"
	cmd "api-gateway/v2/modules/krakend-cobra/v2"
	flexibleconfig "api-gateway/v2/modules/krakend-flexibleconfig/v2"
	viper "api-gateway/v2/modules/krakend-viper/v2"
	"api-gateway/v2/modules/lura/v2/config"
	"log"
	krakend "api-gateway/v2"
	"os"
	"os/signal"
	"syscall"
)

const (
	fcPartials  = "FC_PARTIALS"
	fcTemplates = "FC_TEMPLATES"
	fcSettings  = "FC_SETTINGS"
	fcPath      = "FC_OUT"
	fcEnable    = "FC_ENABLE"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case sig := <-sigs:
			log.Println("Signal intercepted:", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	krakend.RegisterEncoders()

	for key, alias := range aliases {
		config.ExtraConfigAlias[alias] = key
	}

	var cfg config.Parser
	cfg = viper.New()
	if os.Getenv(fcEnable) != "" {
		cfg = flexibleconfig.NewTemplateParser(flexibleconfig.Config{
			Parser:    cfg,
			Partials:  os.Getenv(fcPartials),
			Settings:  os.Getenv(fcSettings),
			Path:      os.Getenv(fcPath),
			Templates: os.Getenv(fcTemplates),
		})
	}

	cmd.Execute(cfg, krakend.NewExecutor(ctx))
}

var aliases = map[string]string{
	"github_com/davron112/lure/transport/http/server/handler":  "plugin/http-server",
	"api-gateway/v2/modules/lura/transport/http/client/executor": "plugin/http-client",
	"api-gateway/v2/modules/lura/proxy/plugin":                   "plugin/req-resp-modifier",
	"api-gateway/v2/modules/lura/proxy":                          "proxy",
	"github_com/davron112/lura/router/gin":                     "router",

	"api-gateway/v2/modules/krakend-httpcache":                "qos/http-cache",
	"api-gateway/v2/modules/krakend-circuitbreaker/gobreaker": "qos/circuit-breaker",

	"api-gateway/v2/modules/krakend-oauth2-clientcredentials": "auth/client-credentials",
	"api-gateway/v2/modules/krakend-jose/validator":           "auth/validator",
	"api-gateway/v2/modules/krakend-jose/signer":              "auth/signer",
	"github_com/davron112/bloomfilter":                      "auth/revoker",

	"github_com/davron112/krakend-botdetector": "security/bot-detector",
	"github_com/davron112/krakend-httpsecure":  "security/http",
	"github_com/davron112/krakend-cors":        "security/cors",

	"api-gateway/v2/modules/krakend-cel":        "validation/cel",
	"api-gateway/v2/modules/krakend-jsonschema": "validation/json-schema",

	"api-gateway/v2/modules/krakend-amqp/agent": "async/amqp",

	"api-gateway/v2/modules/krakend-amqp/consume":                  "backend/amqp/consumer",
	"api-gateway/v2/modules/krakend-amqp/produce":                  "backend/amqp/producer",
	"api-gateway/v2/modules/krakend-lambda":                        "backend/lambda",
	"api-gateway/v2/modules/krakend-pubsub/publisher":              "backend/pubsub/publisher",
	"api-gateway/v2/modules/krakend-pubsub/subscriber":             "backend/pubsub/subscriber",
	"api-gateway/v2/modules/krakend/transport/http/client/graphql": "backend/graphql",
	"api-gateway/v2/modules/grpc":                                  "backend/grpc",
	"api-gateway/v2/modules/krakend/http":                          "backend/http",

	"github_com/davron112/krakend-gelf":       "telemetry/gelf",
	"github_com/davron112/krakend-gologging":  "telemetry/logging",
	"api-gateway/v2/modules/krakend-logstash":   "telemetry/logstash",
	"github_com/davron112/krakend-metrics":    "telemetry/metrics",
	"github_com/davron112/krakend-influx":     "telemetry/influx",
	"github_com/davron112/krakend-opencensus": "telemetry/opencensus",

	"api-gateway/v2/modules/krakend-lua/router":        "modifier/lua-endpoint",
	"api-gateway/v2/modules/krakend-lua/proxy":         "modifier/lua-proxy",
	"api-gateway/v2/modules/krakend-lua/proxy/backend": "modifier/lua-backend",
	"api-gateway/v2/modules/krakend-martian":           "modifier/martian",
}
