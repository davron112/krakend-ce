package audit

import (
	"time"

	bf "api-gateway/v2/modules/bloomfilter/v2/krakend"
	botdetector "api-gateway/v2/modules/krakend-botdetector/v2/krakend"
	opencensus "api-gateway/v2/modules/krakend-opencensus/v2"
	juju "api-gateway/v2/modules/krakend-ratelimit/v2/juju/router"
	rss "api-gateway/v2/modules/krakend-rss/v2"
	xml "api-gateway/v2/modules/krakend-xml/v2"
	"api-gateway/v2/modules/lura/v2/config"
	"api-gateway/v2/modules/lura/v2/encoding"
	"api-gateway/v2/modules/lura/v2/proxy"
	"api-gateway/v2/modules/lura/v2/proxy/plugin"
	router "api-gateway/v2/modules/lura/v2/router/gin"
	client "api-gateway/v2/modules/lura/v2/transport/http/client/plugin"
	server "api-gateway/v2/modules/lura/v2/transport/http/server/plugin"
)

// Parse creates a Service capturing the details of the received configuration
func Parse(cfg *config.ServiceConfig) Service {
	v1 := 0

	if cfg.Plugin != nil {
		v1 = addBit(v1, ServicePlugin)
	}

	if cfg.SequentialStart {
		v1 = addBit(v1, ServiceSequentialStart)
	}

	if cfg.Debug {
		v1 = addBit(v1, ServiceDebug)
	}

	if cfg.AllowInsecureConnections {
		v1 = addBit(v1, ServiceAllowInsecureConnections)
	}

	if cfg.DisableStrictREST {
		v1 = addBit(v1, ServiceDisableStrictREST)
	}

	if cfg.TLS != nil {
		v1 = addBit(v1, ServiceHasTLS)
		if !cfg.TLS.IsDisabled {
			v1 = addBit(v1, ServiceTLSEnabled)
		}
		if cfg.TLS.EnableMTLS {
			v1 = addBit(v1, ServiceTLSEnableMTLS)
		}
		if cfg.TLS.DisableSystemCaPool {
			v1 = addBit(v1, ServiceTLSDisableSystemCaPool)
		}
		if len(cfg.TLS.CaCerts) > 0 {
			v1 = addBit(v1, ServiceTLSCaCerts)
		}
	}

	return Service{
		Details:    []int{v1},
		Agents:     parseAsyncAgents(cfg.AsyncAgents),
		Endpoints:  parseEndpoints(cfg.Endpoints),
		Components: parseComponents(cfg.ExtraConfig),
	}
}

func parseAsyncAgents(as []*config.AsyncAgent) []Agent {
	agents := []Agent{}
	for _, a := range as {
		agent := Agent{
			Details: []int{
				parseEncoding(a.Encoding),
				a.Consumer.Workers,
				a.Connection.MaxRetries,
				int(a.Consumer.Timeout / time.Millisecond),
			},
			Backends:   parseBackends(a.Backend),
			Components: parseComponents(a.ExtraConfig),
		}

		agents = append(agents, agent)
	}
	return agents
}

func parseEndpoints(es []*config.EndpointConfig) []Endpoint {
	endpoints := []Endpoint{}
	for _, e := range es {
		endpoint := Endpoint{
			Details: []int{
				parseEncoding(e.OutputEncoding),
				len(e.QueryString),
				len(e.HeadersToPass),
				int(e.Timeout / time.Millisecond),
			},
			Backends:   parseBackends(e.Backend),
			Components: parseComponents(e.ExtraConfig),
		}

		endpoints = append(endpoints, endpoint)
	}
	return endpoints
}

func parseEncoding(enc string) int {
	switch enc {
	case encoding.NOOP:
		return addBit(0, EncodingNOOP)
	case encoding.JSON:
		return addBit(0, EncodingJSON)
	case encoding.SAFE_JSON:
		return addBit(0, EncodingSAFEJSON)
	case encoding.STRING:
		return addBit(0, EncodingSTRING)
	case rss.Name:
		return addBit(0, EncodingRSS)
	case xml.Name:
		return addBit(0, EncodingXML)
	default:
		return addBit(0, EncodingOther)
	}
}

func parseBackends(bs []*config.Backend) []Backend {
	backends := []Backend{}
	for _, b := range bs {
		v1 := parseEncoding(b.Encoding)
		if len(b.AllowList) > 0 {
			v1 = addBit(v1, BackendAllow)
		}
		if len(b.DenyList) > 0 {
			v1 = addBit(v1, BackendDeny)
		}
		if len(b.Mapping) > 0 {
			v1 = addBit(v1, BackendMapping)
		}
		if b.Group != "" {
			v1 = addBit(v1, BackendGroup)
		}
		if b.Target != "" {
			v1 = addBit(v1, BackendTarget)
		}
		if b.IsCollection {
			v1 = addBit(v1, BackendIsCollection)
		}
		backend := Backend{
			Details:    []int{v1},
			Components: parseComponents(b.ExtraConfig),
		}

		backends = append(backends, backend)
	}
	return backends
}

func parseComponents(cfg config.ExtraConfig) Component {
	components := Component{}
	for c, v := range cfg {
		switch c {
		case server.Namespace:
			cfg, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			if n, ok := cfg["name"].(string); ok {
				components[c] = []int{addBit(0, parseServerPlugin(n))}
				continue
			}

			if ns, ok := cfg["name"].([]interface{}); ok {
				vs := 0
				for _, raw := range ns {
					n, ok := raw.(string)
					if !ok {
						continue
					}
					vs = addBit(vs, parseServerPlugin(n))
				}
				components[c] = []int{vs}
				continue
			}

		case client.Namespace:
			cfg, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			n, ok := cfg["name"].(string)
			if !ok {
				continue
			}
			components[c] = []int{parseClientPlugin(n)}

		case plugin.Namespace:
			cfg, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			ns, ok := cfg["name"].([]interface{})
			if !ok {
				continue
			}
			vs := 0
			for _, raw := range ns {
				n, ok := raw.(string)
				if !ok {
					continue
				}
				vs = addBit(vs, parseRespReqPlugin(n))
			}
			components[c] = []int{vs}

		case proxy.Namespace:
			cfg, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			components[c] = []int{parseProxy(cfg)}

		case router.Namespace:
			cfg, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			components[c] = []int{parseRouter(cfg)}

		case bf.Namespace:
			cfg, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			res := make([]int, 2)
			if hn, ok := cfg["hash_name"].(string); ok && hn == "optimal" {
				res[0] = 1
			}
			if ks, ok := cfg["token_keys"].([]interface{}); ok {
				res[1] = len(ks)
			}
			components[c] = res

		case botdetector.Namespace:
			cfg, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			res := make([]int, 4)
			if ks, ok := cfg["allow"].([]interface{}); ok {
				res[0] = len(ks)
			}
			if ks, ok := cfg["deny"].([]interface{}); ok {
				res[1] = len(ks)
			}
			if ks, ok := cfg["patterns"].([]interface{}); ok {
				res[2] = len(ks)
			}
			if s, ok := cfg["cache_size"].(float64); ok {
				res[3] = int(s)
			}
			components[c] = res

		case opencensus.Namespace:
			cfg, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			exp, ok := cfg["exporters"].(map[string]interface{})
			if !ok {
				continue
			}

			v1 := 0
			if _, ok := exp["logger"]; ok {
				v1 = 1
			}
			if _, ok := exp["zipkin"]; ok {
				v1 += 2
			}
			if _, ok := exp["jaeger"]; ok {
				v1 += 4
			}
			if _, ok := exp["influxdb"]; ok {
				v1 += 8
			}
			if _, ok := exp["prometheus"]; ok {
				v1 += 16
			}
			if _, ok := exp["xray"]; ok {
				v1 += 32
			}
			if _, ok := exp["stackdriver"]; ok {
				v1 += 64
			}
			if _, ok := exp["datadog"]; ok {
				v1 += 128
			}
			if _, ok := exp["ocagent"]; ok {
				v1 += 256
			}

			components[c] = []int{v1}

		case juju.Namespace:
			cfg, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			v1 := 0
			if vs, ok := cfg["max_rate"].(float64); ok && vs > 0 {
				v1 = 1
			}
			if vs, ok := cfg["client_max_rate"].(float64); ok && vs > 0 {
				v1 += 2
			}
			if vs, ok := cfg["strategy"].(string); ok {
				switch vs {
				case "ip":
					v1 += 4
				case "header":
					v1 += 8
				}
			}

			components[c] = []int{v1}

		default:
			components[c] = []int{}
		}
	}
	return components
}

func parseRouter(cfg config.ExtraConfig) int {
	res := 0
	v, ok := cfg["error_body"].(bool)
	if ok && v {
		res = addBit(res, RouterErrorBody)
	}

	v, ok = cfg["disable_health"].(bool)
	if ok && v {
		res = addBit(res, RouterDisableHealth)
	}

	v, ok = cfg["disable_access_log"].(bool)
	if ok && v {
		res = addBit(res, RouterDisableAccessLog)
	}

	if _, ok := cfg["health_path"]; ok {
		res = addBit(res, RouterHealthPath)
	}

	v, ok = cfg["return_error_msg"].(bool)
	if ok && v {
		res = addBit(res, RouterErrorMsg)
	}

	v, ok = cfg["disable_redirect_trailing_slash"].(bool)
	if ok && v {
		res = addBit(res, RouterDisableRedirectTrailingSlash)
	}

	v, ok = cfg["disable_redirect_fixed_path"].(bool)
	if ok && v {
		res = addBit(res, RouterDisableRedirectFixedPath)
	}

	v, ok = cfg["remove_extra_slash"].(bool)
	if ok && v {
		res = addBit(res, RouterExtraSlash)
	}

	v, ok = cfg["disable_handle_method_not_allowed"].(bool)
	if ok && v {
		res = addBit(res, RouterHandleMethodNotAllowed)
	}

	v, ok = cfg["disable_path_decoding"].(bool)
	if ok && v {
		res = addBit(res, RouterPathDecoding)
	}

	v, ok = cfg["auto_options"].(bool)
	if ok && v {
		res = addBit(res, RouterAutoOptions)
	}

	v, ok = cfg["forwarded_by_client_ip"].(bool)
	if ok && v {
		res = addBit(res, RouterForwardedByClientIp)
	}

	vs, ok := cfg["remote_ip_headers"].([]interface{})
	if ok && len(vs) > 0 {
		res = addBit(res, RouterRemoteIpHeaders)
	}

	vs, ok = cfg["trusted_proxies"].([]interface{})
	if ok && len(vs) > 0 {
		res = addBit(res, RouterTrustedProxies)
	}

	v, ok = cfg["app_engine"].(bool)
	if ok && v {
		res = addBit(res, RouterAppEngine)
	}

	if v, ok := cfg["max_multipart_memory"].(float64); ok && v > 0 {
		res = addBit(res, RouterMaxMultipartMemory)
	}

	vs, ok = cfg["logger_skip_paths"].([]interface{})
	if ok && len(vs) > 0 {
		res = addBit(res, RouterLoggerSkipPaths)
	}

	v, ok = cfg["hide_version_header"].(bool)
	if ok && v {
		res = addBit(res, RouterHideVersionHeader)
	}

	return res
}

func parseProxy(cfg config.ExtraConfig) int {
	res := 0
	v, ok := cfg["sequential"].(bool)
	if ok && v {
		res = addBit(res, 0)
	}

	if _, ok := cfg["flatmap_filter"]; ok {
		res = addBit(res, 1)
	}

	v, ok = cfg["shadow"].(bool)
	if ok && v {
		res = addBit(res, 2)
	}

	if _, ok := cfg["combiner"]; ok {
		res = addBit(res, 3)
	}

	if _, ok := cfg["static"]; ok {
		res = addBit(res, 4)
	}
	return res
}

func parseServerPlugin(name string) int {
	switch name {
	case "static-filesystem":
		return 1
	case "basic-auth":
		return 2
	case "geoip":
		return 3
	case "redis-ratelimit":
		return 4
	case "url-rewrite":
		return 5
	case "virtualhost":
		return 6
	case "wildcard":
		return 7
	case "ip-filter":
		return 8
	case "jwk-aggregator":
		return 9
	}
	return 0
}

func parseClientPlugin(name string) int {
	switch name {
	case "no-redirect":
		return 1
	case "http-logger":
		return 2
	case "static-filesystem":
		return 3
	case "http-proxy":
		return 4
	}
	return 0
}

func parseRespReqPlugin(name string) int {
	switch name {
	case "response-schema-validator":
		return 1
	case "content-replacer":
		return 2
	}
	return 0
}

func addBit(x, y int) int {
	return x | (1 << y)
}
