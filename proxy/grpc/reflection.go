package grpc

import (
	"context"
	"api-gateway/v2/modules/lura/v2/logging"
	"api-gateway/v2/modules/lura/v2/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/pkg/errors"
)

type Reflector struct {
	c      *grpcreflect.Client
	logger logging.Logger
}

// NewReflector creates a new Reflector from the reflection client and a lura logger.
func NewReflector(client *grpcreflect.Client, logger logging.Logger) *Reflector {
	return &Reflector{
		c:      client,
		logger: logger,
	}
}

// CreateInvocation creates a MethodInvocation by performing reflection
func (r *Reflector) CreateInvocation(ctx context.Context, serviceName, methodName string, input []byte) (*MethodInvocation, error) {
	serviceDesc, err := r.c.ResolveService(serviceName)
	if err != nil {
		r.logger.Error("Failed to resolve service", map[string]interface{}{
			"service": serviceName,
			"error":   err.Error(),
		})
		return nil, err
	}

	methodDesc := serviceDesc.FindMethodByName(methodName)
	if methodDesc == nil {
		r.logger.Error("Method not found in service", map[string]interface{}{
			"method":  methodName,
			"service": serviceName,
		})
		return nil, errors.New("method not found upstream")
	}

	inputMessage := dynamic.NewMessage(methodDesc.GetInputType())
	if methodDesc.GetInputType().GetName() == "Empty" && len(input) > 0 {
		r.logger.Warning("Non-empty body received for an operation expecting an Empty message type.", map[string]interface{}{
			"service": serviceName,
			"method":  methodName,
		})
		return nil, utils.NewHTTPError(500, "Unexpected non-empty body for Empty type")
	} else if len(input) > 0 {
		if err := inputMessage.UnmarshalJSON(input); err != nil {
			r.logger.Error("Failed to unmarshal request into dynamic message", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, utils.NewHTTPError(500, "Failed to unmarshal request")
		}
	}

	return &MethodInvocation{
		MethodDescriptor: methodDesc,
		Message:          inputMessage,
	}, nil
}

// MethodInvocation contains a method and a message used to invoke an RPC
type MethodInvocation struct {
	*desc.MethodDescriptor
	*dynamic.Message
}
