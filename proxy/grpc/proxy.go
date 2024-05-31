package grpc

import (
	"context"
	"encoding/json"
	"api-gateway/v2/modules/lura/v2/logging"
	"github.com/davron112/lura/v2/utils"
	"github.com/golang/protobuf/jsonpb"
	"net/url"

	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

type Proxy struct {
	cc        *grpc.ClientConn
	reflector *Reflector
	stub      grpcdynamic.Stub
}

// NewProxy creates a new client
func NewProxy() *Proxy {
	return &Proxy{}
}

// Connect opens a connection to target.
func (p *Proxy) Connect(ctx context.Context, target *url.URL, logger logging.Logger) error {
	cc, err := grpc.DialContext(ctx, target.String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	p.cc = cc
	rc := grpcreflect.NewClientV1Alpha(ctx, rpb.NewServerReflectionClient(p.cc))
	p.reflector = NewReflector(rc, logger)
	p.stub = grpcdynamic.NewStub(p.cc)
	return err
}

// Call performs the gRPC call after doing reflection to obtain type information.
func (p *Proxy) Call(ctx context.Context, serviceName, methodName string, message []byte, md *metadata.MD) ([]byte, error) {
	invocation, err := p.reflector.CreateInvocation(ctx, serviceName, methodName, message)
	if err != nil {
		return nil, err
	}

	output, err := p.stub.InvokeRpc(ctx, invocation.MethodDescriptor, invocation.Message, grpc.Header(md))
	if err != nil {
		stat := status.Convert(err)
		grpcError := utils.NewHTTPError(int(stat.Code()), stat.Message())
		grpcErrorJson, _ := json.Marshal(grpcError)
		return grpcErrorJson, grpcError
	}

	outputMessage := dynamic.NewMessage(invocation.MethodDescriptor.GetOutputType())
	err = outputMessage.ConvertFrom(output)
	if err != nil {
		return nil, errors.Wrap(err, "response from backend could not be converted internally")
	}

	m, err := outputMessage.MarshalJSONPB(&jsonpb.Marshaler{EnumsAsInts: true, EmitDefaults: true})
	if err != nil {
		return nil, err
	}
	return m, err
}
