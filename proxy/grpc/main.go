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
	Namespace      = "github.com/davron112/grpc"
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
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// IsGrpcMethod checks if the given backend configuration is designated for gRPC.
func IsGrpcMethod(remote *config.Backend) bool {
	grpcKeys := []string{"backend/grpc", "github.com/davron112/grpc"}

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
		if !IsGrpcMethod(remote) {
			return next
		}

		getNextHost := RoundRobinHostSelector(remote.Host)

		return func(requestCtx context.Context, req *proxy.Request) (*proxy.Response, error) {
			timeoutDuration := 3000 * time.Millisecond
			ctx, cancel := context.WithTimeout(requestCtx, timeoutDuration)
			defer cancel()

			hostURL, err := url.Parse(trimPrefixes(getNextHost()))
			if err != nil {
				logger.Error("URL parsing error: ", err)
				return createErrorResponse(500, "URL parsing error"), nil
			}

			grpcProxy := NewProxy()
			if err := grpcProxy.Connect(ctx, hostURL, logger); err != nil {
				logger.Error("Failed to connect: ", err)
				return createErrorResponse(500, "Failed to connect"), nil
			}

			bodyBytes, err := io.ReadAll(req.Body)
			req.Body.Close()
			if err != nil {
				logger.Error("Failed to read request body: ", err)
				return createErrorResponse(500, "Failed to read request body"), nil
			}

			serviceName, methodName, err := parseURLPattern(remote.URLPattern)
			if err != nil {
				logger.Error(err.Error())
				return createErrorResponse(400, err.Error()), nil
			}

			md := metadata.New(nil)
			responseBytes, err := grpcProxy.Call(ctx, serviceName, methodName, bodyBytes, &md)
			if err != nil {
				logger.Error("gRPC call failed: ", err)
				return createErrorResponse(500, err.Error()), nil
			}

			var responseData map[string]interface{}
			if err := json.Unmarshal(responseBytes, &responseData); err != nil {
				logger.Error("Failed to unmarshal JSON data into map: ", err)
				return createErrorResponse(500, "Failed to unmarshal JSON data into map"), nil
			}

			return &proxy.Response{
				Data:       responseData,
				IsComplete: true,
			}, nil
		}
	}
}

func parseURLPattern(urlPattern string) (serviceName, methodName string, err error) {
	if !strings.HasPrefix(urlPattern, "/") {
		return "", "", errors.New("URL pattern must start with a '/'")
	}

	parts := strings.Split(urlPattern, "/")
	if len(parts) != 3 {
		return "", "", errors.New("Invalid URL pattern format. Expected format: /ServiceName/MethodName")
	}

	serviceName = parts[1]
	methodName = parts[2]
	return serviceName, methodName, nil
}

func createErrorResponse(statusCode int, message string) *proxy.Response {
	responseData := map[string]interface{}{
		"msg":    message,
	}
	return &proxy.Response{
		Data:       responseData,
		IsComplete: false,
		Metadata:   proxy.Metadata{StatusCode: statusCode},
	}
}

