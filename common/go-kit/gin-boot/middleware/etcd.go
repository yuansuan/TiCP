package middleware

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"

	"go.etcd.io/etcd/client/pkg/v3/transport"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (mw *Middleware) initEtcd() {
	conf := &mw.conf.App.Middleware.Etcd
	if !conf.Startup {
		return
	}

	clientConf := clientv3.Config{
		Endpoints:   mw.conf.App.Middleware.Etcd.Endpoints,
		DialTimeout: time.Second * 20,
	}
	if conf.TLS {
		tlsinfo := transport.TLSInfo{
			CertFile:      conf.CertFile,
			KeyFile:       conf.KeyFile,
			TrustedCAFile: conf.CAFile,
		}

		tls, err := tlsinfo.ClientConfig()
		if err != nil {
			panic(err)
		}

		clientConf.TLS = tls
	}

	client, err := clientv3.New(clientConf)
	if err != nil {
		util.ChkErr(err)
	}
	mw.etcd.client = client
	mw.etcd.polling = conf.Polling
}
