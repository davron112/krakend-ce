package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"api-gateway/v2/modules/lura/v2/config"
	"api-gateway/v2/modules/lura/v2/logging"
	"api-gateway/v2/modules/lura/v2/proxy"
	"api-gateway/v2/modules/lura/v2/utils"
	"google.golang.org/grpc/metadata"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	Namespace      = "api-gateway/v2/modules/grpc"
	EncodingPrefix = "grpc-"
)

func trimPrefixes(s string) string {
	re, err := regexp.Compile(`^(?:http://|https://|gpc://)`)
	if err != nil {
		fmt.Println("Regex compile error:", err)
		return s
	}
	return re.ReplaceAllString(s, "")
}

func RoundRobinHostSelector(hosts []string) func() string {
	var next int
	return func() string {
		host := hosts[next%len(hosts)]
		next++
		return host
	}
}

// ErrorResponse to encapsulate an error response
type ErrorResponse struct {
	StatusCode int
	Message    string
}

// IsGrpcMethod checks if the given backend configuration is designated for gRPC.
func IsGrpcMethod(remote *config.Backend) bool {
	// List of potential keys that signify a gRPC backend.
	grpcKeys := []string{"backend/grpc", "api-gateway/v2/modules/grpc"}

	// Check if any of the keys exist in the ExtraConfig map.
	for _, key := range grpcKeys {
		if _, exists := remote.ExtraConfig[key]; exists {
			return true
		}
	}

	return false
}

func NewGrpcBackendFactory(logger logging.Logger, bf proxy.BackendFactory) proxy.BackendFactory {

	return func(remote *config.Backend) proxy.Proxy {
		next := bf(remote)
		ok := IsGrpcMethod(remote)
		if !ok {
			return next
		}

		getNextHost := RoundRobinHostSelector(remote.Host)

		return func(requestCtx context.Context, req *proxy.Request) (*proxy.Response, error) {
			timeoutDuration := 3000 * time.Millisecond
			ctx, cancel := context.WithTimeout(requestCtx, timeoutDuration)
			defer cancel()

			// Prepare host
			hostURL, err := url.Parse(trimPrefixes(getNextHost()))
			if err != nil {
				logger.Error("URL parsing error: ", err)
				return nil, utils.NewHTTPError(500, "URL parsing error")
			}

			grpcProxy := NewProxy()
			if err := grpcProxy.Connect(ctx, hostURL, logger); err != nil {
				logger.Error("Failed to connect: ", err)
				return nil, utils.NewHTTPError(500, "Failed to connect")
			}

			// Read request data
			bodyBytes, err := io.ReadAll(req.Body)
			req.Body.Close()
			if err != nil {
				logger.Error("Failed to read request body: ", err)
				return nil, utils.NewHTTPError(500, "Failed to read request body")
			}

			// Prepare URLPatter
			serviceName, methodName, err := parseURLPattern(remote.URLPattern)
			if err != nil {
				logger.Error(err.Error())
				return nil, utils.NewHTTPError(400, err.Error())
			}

			// Send request
			fmt.Println(serviceName, methodName, "serviceName, methodName")
			md := metadata.New(nil) // Send metadata here, if needed
			responseBytes, err := grpcProxy.Call(ctx, serviceName, methodName, bodyBytes, &md)
			if err != nil {
				logger.Error("gRPC call failed: ", err)
				return nil, utils.NewHTTPError(500, "gRPC call failed")
			}

			// Prepare json response
			var responseData map[string]interface{}
			if err := json.Unmarshal(responseBytes, &responseData); err != nil {
				logger.Error("Failed to unmarshal JSON data into map: ", err)
				return nil, utils.NewHTTPError(500, "Failed to unmarshal JSON data into map")
			}

			return &proxy.Response{
				Data:       responseData,
				IsComplete: true,
			}, nil
		}
	}
}

func parseURLPattern(urlPattern string) (serviceName, methodName string, err error) {
	// Ensure the URL pattern starts with a slash ("/").
	if !strings.HasPrefix(urlPattern, "/") {
		return "", "", errors.New("URL pattern must start with a '/'")
	}

	// Split the URL pattern by slash ("/") to separate components.
	parts := strings.Split(urlPattern, "/")
	// Expecting three parts: an empty string (due to leading slash), the service name, and the method name.
	if len(parts) != 3 {
		return "", "", errors.New("Invalid URL pattern format. Expected format: /ServiceName/MethodName")
	}

	// Assign and return the service and method names.
	serviceName = parts[1]
	methodName = parts[2]
	return serviceName, methodName, nil
}
