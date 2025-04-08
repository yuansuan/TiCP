package grpc_boot

import (
	"errors"
	"math"
	"net"
	"strconv"
	"sync"
	"time"

	conf_type "github.com/yuansuan/ticp/common/go-kit/gin-boot/conf-type"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/keepalive"
)

var (
	server                           = &ServerType{}
	listener                         = net.Listener(nil)
	confServer                       = &conf_type.GRPCServer{}
	defaultKeepaliveServerParameters = keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(math.MaxInt64),
		MaxConnectionAge:      time.Duration(math.MaxInt64),
		MaxConnectionAgeGrace: time.Duration(math.MaxInt64),
		Time:                  7200 * time.Second,
		Timeout:               20 * time.Second,
	}

	defaultKeepaliveEnforcementPolicy = keepalive.EnforcementPolicy{
		MinTime:             300 * time.Second,
		PermitWithoutStream: false,
	}
)

// ServerType is type of Server
type ServerType struct {
	name     string
	driver   *grpc.Server
	listener net.Listener
}

// InitListener InitListener
func InitListener(addr string) error {
	var err error
	listener, err = net.Listen("tcp", addr)
	return err
}

// InitServer InitServer
func InitServer(conf *conf_type.GRPCServer) {
	confServer = conf
	s := Create("invalid" /*config.Conf.App.Cluster*/, config.Conf.App.Name)
	server = s
	service.RegisterChannelzServiceToServer(s.Driver())
}

var serverLock sync.Mutex

// DefaultServer DefaultServer
func DefaultServer() (*ServerType, error) {
	name := "default"
	serverLock.Lock()
	defer serverLock.Unlock()
	if server != nil {
		return server, nil
	}
	return nil, errors.New("gRPC server is unavailable: " + name)
}

func unaryServerChain(interceptor ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	arr := []grpc.UnaryServerInterceptor{}
	for _, i := range interceptor {
		if i != nil {
			arr = append(arr, i)
		}
	}
	return grpc_middleware.WithUnaryServerChain(arr...)
}

// Create creates grpc server
func Create(cluster, name string) *ServerType {
	s := &ServerType{
		name: name,
		driver: grpc.NewServer(append(
			getServerOption(confServer),
			unaryServerChain(middleware.GrpcInterceptors...),
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				middleware.GrpcStreamInterceptors...,
			)),
		)...),
		listener: listener,
	}

	return s
}

// Driver Driver
func (s *ServerType) Driver() *grpc.Server {
	return s.driver
}

// Listener Listener
func (s *ServerType) Listener() net.Listener {
	return s.listener
}

// Servers is servers?
func Server() *ServerType {
	return server
}

// getServerOption 根据配置获取serverOption
func getServerOption(options *conf_type.GRPCServer) []grpc.ServerOption {
	serverOptions := []grpc.ServerOption{}

	if options.WriteBufferSize != "" {
		writeBufferSize, err := strconv.Atoi(options.WriteBufferSize)
		if err != nil {
			util.ChkErr(errors.New("rpc server write_buffer_size config error"))
		}
		serverOptions = append(serverOptions, grpc.WriteBufferSize(writeBufferSize))
	}

	if options.ReadBufferSize != "" {
		readBufferSize, err := strconv.Atoi(options.ReadBufferSize)
		if err != nil {
			util.ChkErr(errors.New("rpc server read_buffer_size config error"))
		}
		serverOptions = append(serverOptions, grpc.ReadBufferSize(readBufferSize))
	}

	if options.InitialWindowSize > 0 {
		serverOptions = append(serverOptions, grpc.InitialWindowSize(options.InitialWindowSize))
	}

	if options.InitialConnWindowSize > 0 {
		serverOptions = append(serverOptions, grpc.InitialConnWindowSize(options.InitialConnWindowSize))
	}

	if options.KeepaliveMaxConnectionIdle > 0 ||
		options.KeepaliveMaxConnectionAge > 0 ||
		options.KeepaliveMaxConnectionAgeGrace > 0 ||
		options.KeepaliveTime > 0 ||
		options.KeepaliveTimeout > 0 {
		KeepaliveServerParams := defaultKeepaliveServerParameters
		if options.KeepaliveMaxConnectionIdle > 0 {
			KeepaliveServerParams.MaxConnectionIdle = options.KeepaliveMaxConnectionIdle * time.Millisecond
		}
		if options.KeepaliveMaxConnectionAge > 0 {
			KeepaliveServerParams.MaxConnectionAge = options.KeepaliveMaxConnectionAge * time.Millisecond
		}
		if options.KeepaliveMaxConnectionAgeGrace > 0 {
			KeepaliveServerParams.MaxConnectionAgeGrace = options.KeepaliveMaxConnectionAgeGrace * time.Millisecond
		}
		if options.KeepaliveTime > 0 {
			KeepaliveServerParams.Time = options.KeepaliveTime * time.Millisecond
		}
		if options.KeepaliveTimeout > 0 {
			KeepaliveServerParams.Timeout = options.KeepaliveTimeout * time.Millisecond
		}
		serverOptions = append(serverOptions, grpc.KeepaliveParams(KeepaliveServerParams))
	}

	if options.KeepaliveEnforcementPolicyMinTime > 0 || options.KeepaliveEnforcementPolicyPermitWithoutStream != false {
		keepaliveEnforcementPolicy := defaultKeepaliveEnforcementPolicy
		if options.KeepaliveEnforcementPolicyMinTime > 0 {
			keepaliveEnforcementPolicy.MinTime = options.KeepaliveEnforcementPolicyMinTime * time.Millisecond
		}
		if options.KeepaliveEnforcementPolicyPermitWithoutStream == true {
			keepaliveEnforcementPolicy.PermitWithoutStream = options.KeepaliveEnforcementPolicyPermitWithoutStream
		}
		serverOptions = append(serverOptions, grpc.KeepaliveEnforcementPolicy(keepaliveEnforcementPolicy))
	}

	if options.MaxRecvMsgSize != "" {
		maxRecvMsgSize, err := strconv.Atoi(options.MaxRecvMsgSize)
		if err != nil {
			util.ChkErr(errors.New("rpc server max_recv_msg_size config error"))
		}
		serverOptions = append(serverOptions, grpc.MaxRecvMsgSize(maxRecvMsgSize))
	}

	if options.MaxSendMsgSize != "" {
		maxSendMsgSize, err := strconv.Atoi(options.MaxSendMsgSize)
		if err != nil {
			util.ChkErr(errors.New("rpc server max_send_msg_size config error"))
		}
		serverOptions = append(serverOptions, grpc.MaxSendMsgSize(maxSendMsgSize))
	}

	if options.MaxConcurrentStreams > 0 {
		serverOptions = append(serverOptions, grpc.MaxConcurrentStreams(options.MaxConcurrentStreams))
	}

	if options.ConnectionTimeout > 0 {
		serverOptions = append(serverOptions, grpc.ConnectionTimeout(options.ConnectionTimeout*time.Millisecond))
	}

	if options.MaxHeaderListSize > 0 {
		serverOptions = append(serverOptions, grpc.MaxHeaderListSize(options.MaxHeaderListSize))
	}

	return serverOptions
}
