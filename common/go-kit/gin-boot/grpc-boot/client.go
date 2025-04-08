package grpc_boot

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"go.elastic.co/apm/module/apmgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"

	conf_type "github.com/yuansuan/ticp/common/go-kit/gin-boot/conf-type"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
)

const clientInjectTag = "grpc_client_inject"
const defaultClientTimeoutPerRequest = time.Second * 5
const clientPerReqTimeout = "cli_req_tmout"

var passThroughMetadata = []string{middleware.UserKey, middleware.TraceID}

var (
	clients                = map[string]*grpc.ClientConn{}
	confClients            = &conf_type.GRPCClients{}
	defaultKeepAliveParams = keepalive.ClientParameters{
		Time:                time.Duration(math.MaxInt64),
		Timeout:             20 * time.Second,
		PermitWithoutStream: false,
	}

	// SomeServiceClient => func(*grpc.ClientConn) SomeServiceClient
	registered = map[reflect.Type]*reflect.Value{}
)

// InjectAllClient InjectAllClient
func InjectAllClient(v interface{}) {
	log := logging.Default().With("function", "inject_grpc_client")
	if v == nil {
		log.Panic("unable inject client to nil value")
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		log.Panic("unable bind to non-pointer value")
	}

	if rv = rv.Elem(); rv.Kind() != reflect.Struct {
		log.Panic("unsupported type, use struct instead")
	}

	if err := injectClientToStruct(rv); err != nil {
		log.Panic(err)
	}
}

// injectClientToStruct 将对应的 GRPC 客户端注入到结构体中
func injectClientToStruct(v reflect.Value) error {
	rt := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if fv := v.Field(i); fv.CanSet() && fv.CanAddr() {
			ft := fv.Type()

			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
				if fv.IsNil() {
					fv.Set(reflect.New(ft))
				}
				fv = fv.Elem()
			}

			if fv.Kind() == reflect.Struct {
				if err := injectClientToStruct(fv); err != nil {
					return err
				}
				continue
			}

			if fv.Kind() == reflect.Interface {
				if tag := rt.Field(i).Tag.Get(clientInjectTag); len(tag) != 0 {
					if construct, ok := registered[ft]; ok {
						conn, err := GetClient(tag)
						if err != nil {
							return err
						}

						results := construct.Call([]reflect.Value{reflect.ValueOf(conn)})
						fv.Set(results[0]) // it's okay
					}
				}
				continue
			}
		}
	}

	return nil
}

// RegisterClient : RegisterClient("server_kind", func(*grpc.ClientConn) SomeServiceClient)
func RegisterClient(serverKind string, newFunc interface{}) {
	logger := logging.Default().With("error_when_RegisterClient", serverKind)

	tNewFunc := reflect.TypeOf(newFunc)
	vNewFunc := reflect.ValueOf(newFunc)
	if tNewFunc.Kind() != reflect.Func {
		logger.Fatal("newFunc should a Func")
	}
	if tNewFunc.NumIn() != 1 {
		logger.Fatal("newFunc should have/only have one In parameter")
	}
	if tNewFunc.NumOut() != 1 {
		logger.Fatal("newFunc should have/only have one Out result")
	}

	clientType := tNewFunc.Out(0)
	registered[clientType] = &vNewFunc
}

// InitClient InitClient
func InitClient(conf *conf_type.GRPCClients) {
	confClients = conf
}

// DefaultClient DefaultClient
func DefaultClient() (*grpc.ClientConn, error) {
	return GetClient("default")
}

// GetClient GetClient
func GetClient(name string) (*grpc.ClientConn, error) {
	if _, ok := (*confClients)[name]; !ok {
		return nil, fmt.Errorf("not defined in config: grpc client: %v", name)
	}
	cfg := (*confClients)[name]
	dialOptions := GetDialOptions(cfg)
	return NewClient(cfg.Addr, cfg.TimeoutPerRequest.Duration, dialOptions...)
}

// NewClient NewClient
func NewClient(addr string, timeout time.Duration, dialOptions ...grpc.DialOption) (*grpc.ClientConn, error) {
	return dial(addr, timeout, dialOptions...)
}

func copyRequestContextInterceptor(copyFields []string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		outgoingMD, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			outgoingMD = metadata.New(map[string]string{})
		}

		if incomingMD, ok := metadata.FromIncomingContext(ctx); ok {
			for _, field := range copyFields {

				// already set by microservice, skip
				if len(outgoingMD.Get(field)) > 0 {
					continue
				}

				// not set by incomming context, skip
				if len(incomingMD.Get(field)) == 0 {
					continue
				}

				outgoingMD.Set(field, incomingMD.Get(field)...)
			}
			ctx = metadata.NewOutgoingContext(ctx, outgoingMD)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func requestTimeoutClientInterceptor(timeoutPerRequest time.Duration) grpc.UnaryClientInterceptor {
	if timeoutPerRequest == 0 {
		timeoutPerRequest = defaultClientTimeoutPerRequest
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctxTimeout, canceler := context.WithTimeout(ctx, timeoutPerRequest)
		ctxTmout := util.AppendToOutgoingContext(ctxTimeout, clientPerReqTimeout, timeoutPerRequest.String())
		defer canceler()
		return invoker(ctxTmout, method, req, reply, cc, opts...)
	}
}

func dial(addr string, timeoutPerRequest time.Duration, dialOptions ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, append(
		dialOptions,
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
			apmgrpc.NewUnaryClientInterceptor(),
			grpc_prometheus.UnaryClientInterceptor,
			middleware.GRPCEnvClientInterceptor(),
			requestTimeoutClientInterceptor(timeoutPerRequest),
			copyRequestContextInterceptor(passThroughMetadata),
		),
		grpc.WithChainStreamInterceptor(
			otelgrpc.StreamClientInterceptor(),
		),
	)...)
}

// GetDialOptions 根据配置获取
func GetDialOptions(options *conf_type.GRPCClient) []grpc.DialOption {
	var dialOptions []grpc.DialOption

	if options.WithWriteBufferSize != "" {
		writeBufferSize, err := strconv.Atoi(options.WithWriteBufferSize)
		if err != nil {
			util.ChkErr(errors.New("rpc with_write_buffer_size config error"))
		}
		dialOptions = append(dialOptions, grpc.WithWriteBufferSize(writeBufferSize))
	}

	if options.WithReadBufferSize != "" {
		readBufferSize, err := strconv.Atoi(options.WithReadBufferSize)
		if err != nil {
			util.ChkErr(errors.New("rpc with_read_buffer_size config error"))
		}
		dialOptions = append(dialOptions, grpc.WithReadBufferSize(readBufferSize))
	}

	if options.WithInitialWindowSize > 0 {
		dialOptions = append(dialOptions, grpc.WithInitialWindowSize(options.WithInitialWindowSize))
	}

	if options.WithInitialConnWindowSize > 0 {
		dialOptions = append(dialOptions, grpc.WithInitialConnWindowSize(options.WithInitialConnWindowSize))
	}

	if options.WithMaxMsgSize > 0 {
		dialOptions = append(dialOptions, grpc.WithMaxMsgSize(options.WithMaxMsgSize))
	}

	if options.WithBackoffMaxDelay > 0 {
		dialOptions = append(dialOptions, grpc.WithBackoffMaxDelay(options.WithBackoffMaxDelay*time.Millisecond))
	}

	if options.WithBackoffMaxDelay > 0 {
		dialOptions = append(dialOptions, grpc.WithBackoffMaxDelay(options.WithBackoffMaxDelay*time.Millisecond))
	}

	if options.WithBlock == true {
		dialOptions = append(dialOptions, grpc.WithBlock())
	}

	if options.WithInsecure == true {
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	if options.WithTimeout > 0 {
		dialOptions = append(dialOptions, grpc.WithTimeout(options.WithTimeout*time.Millisecond))
	}

	if options.FailOnNonTempDialError == true {
		dialOptions = append(dialOptions, grpc.FailOnNonTempDialError(options.FailOnNonTempDialError))
	}

	if options.WithUserAgent != "" {
		dialOptions = append(dialOptions, grpc.WithUserAgent(options.WithUserAgent))
	}

	if options.WithAuthority != "" {
		dialOptions = append(dialOptions, grpc.WithAuthority(options.WithAuthority))
	}

	if options.WithDisableServiceConfig == true {
		dialOptions = append(dialOptions, grpc.WithDisableServiceConfig())
	}

	if options.WithDisableRetry == true {
		dialOptions = append(dialOptions, grpc.WithDisableRetry())
	}

	if options.WithMaxHeaderListSize > 0 {
		dialOptions = append(dialOptions, grpc.WithMaxHeaderListSize(options.WithMaxHeaderListSize))
	}

	// keepalive ClientParameters
	if options.KeepaliveTime > 0 || options.KeepaliveTimeout > 0 || options.KeepalivePermitWithoutStream != false {
		keepaliveParams := defaultKeepAliveParams
		if options.KeepaliveTime > 0 {
			keepaliveParams.Time = options.KeepaliveTime * time.Millisecond
		}
		if options.KeepaliveTimeout > 0 {
			keepaliveParams.Timeout = options.KeepaliveTimeout * time.Millisecond
		}
		if options.KeepalivePermitWithoutStream != false {
			keepaliveParams.PermitWithoutStream = options.KeepalivePermitWithoutStream
		}
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(keepaliveParams))
	}

	return dialOptions
}
