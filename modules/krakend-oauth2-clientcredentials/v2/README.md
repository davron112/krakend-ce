# KrakenD oauth2 client credentials

A http client for the [KrakenD](https://api-gateway/v2/modules/krakend) framework supporting the oauth2 client credentials workflow.

## How to use it?

This package exposes single factory capable to create a instances of the `proxy.HTTPClientFactory` interface embedding a http client supporting the oauth2 client credentials workflow

	import 	(
		"context"
		"net/http"
		"api-gateway/v2/modules/lura/v2/config"
		"api-gateway/v2/modules/lura/v2/proxy"
		"api-gateway/v2/modules/krakend-oauth2-clientcredentials/v2"
	)

	requestExecutorFactory := func(cfg *config.Backend) proxy.HTTPRequestExecutor {
		clientFactory := oauth2client.NewHTTPClient(cfg)
		return func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return clientFactory(ctx).Do(req.WithContext(ctx))
		}
	}

You can create your own `proxy.HTTPRequestExecutor` and inject it into your `BackendFactory`
