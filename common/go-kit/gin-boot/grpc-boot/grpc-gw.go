package grpc_boot

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strconv"

	boothttp "github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
)

// const strings
const (
	HTTPStatusCodeHeader = "x-http-code"
	HTTPMethodAny        = "any"
)

// GatewayRegisterHandler ...
type GatewayRegisterHandler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

var gatewayHandlers []GatewayRegisterHandler

// RegisterGatewayHandlers 注册grpc client
func RegisterGatewayHandlers(handlers ...GatewayRegisterHandler) {
	gatewayHandlers = append(gatewayHandlers, handlers...)
}

// GatewayURL ...
type GatewayURL struct {
	Method       string
	RelativePath string
}

// with as prefix
var gatewayUrls = []*GatewayURL{
	{Method: "any", RelativePath: "/gapi/*any"},
}

// RegisterGatewayUrls 需要托管的url
func RegisterGatewayUrls(urls ...*GatewayURL) {
	gatewayUrls = append(gatewayUrls, urls...)
}

var cookies = []string{"ys-session-id"}

// RegisterCookieNames 需要转入metadata的cookie
func RegisterCookieNames(names ...string) {
	cookies = append(cookies, names...)
}

var serveMuxOptions []runtime.ServeMuxOption

// RegisterServeMuxOptions 自定义ServeMuxOption
func RegisterServeMuxOptions(options ...runtime.ServeMuxOption) {
	serveMuxOptions = append(serveMuxOptions, options...)
}

func init() {
	serveMuxOptions = []runtime.ServeMuxOption{
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			md := metadata.MD{}
			for _, name := range cookies {
				cookie, err := req.Cookie(name)
				if err == nil {
					md.Set(name, cookie.Value)
				}
			}
			return md
		}),
		runtime.WithOutgoingHeaderMatcher(func(header string) (string, bool) {
			switch header {
			case "set-cookie":
				return header, true
			}
			return "", false
		}),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
			if md, ok := runtime.ServerMetadataFromContext(ctx); ok {
				if vals := md.HeaderMD.Get(HTTPStatusCodeHeader); len(vals) > 0 {
					code, e := strconv.Atoi(vals[0])
					if e == nil {
						err = &runtime.HTTPStatusError{
							HTTPStatus: code,
							Err:        err,
						}
					}
				}
			}
			runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, r, err)
		}),
	}
}

// InitGrpcGateway ...
func InitGrpcGateway(drv *boothttp.Driver) {
	mux := runtime.NewServeMux(serveMuxOptions...)

	addr := config.Conf.App.Middleware.GRPC.Server.Default.Addr
	for _, f := range gatewayHandlers {
		_ = f(context.TODO(), mux, addr, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	}

	for _, v := range gatewayUrls {
		if v.Method == HTTPMethodAny {
			drv.Any(v.RelativePath, gin.WrapF(mux.ServeHTTP))
		} else {
			drv.Handle(v.Method, v.RelativePath, gin.WrapF(mux.ServeHTTP))
		}
	}
}
