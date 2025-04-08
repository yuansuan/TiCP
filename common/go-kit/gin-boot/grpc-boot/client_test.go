package grpc_boot

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	conf_type "github.com/yuansuan/ticp/common/go-kit/gin-boot/conf-type"
)

type DummyClient struct {
	Name string
}

func (d *DummyClient) Action() (string, error) { return d.Name, nil }

type AppleClient interface{ Action() (string, error) }

func NewAppleClient(grpc.ClientConnInterface) AppleClient { return &DummyClient{Name: "apple"} }

type BananaClient interface{ Action() (string, error) }

func NewBananaClient(grpc.ClientConnInterface) BananaClient { return &DummyClient{Name: "banana"} }

const (
	MockServerAddr = "localhost:23456"
)

func init() {
	l, err := net.Listen("tcp", MockServerAddr)
	if err != nil {
		panic(err)
	}

	go func() { _ = grpc.NewServer().Serve(l) }()
}

func TestInjectAllClient(t *testing.T) {
	RegisterClient("_", NewAppleClient)
	RegisterClient("_", NewBananaClient)
	InitClient(&conf_type.GRPCClients{
		"apple": &conf_type.GRPCClient{
			Addr:         "localhost:23456",
			WithInsecure: true,
		},
		"banana": &conf_type.GRPCClient{
			Addr:         "localhost:23456",
			WithInsecure: true,
		},
	})

	t.Run("normal", func(t *testing.T) {
		var v struct {
			Apple  AppleClient  `grpc_client_inject:"apple"`
			Banana BananaClient `grpc_client_inject:"banana"`
		}

		assert.NotPanics(t, func() { InjectAllClient(&v) })
		assert.NotNil(t, v.Apple)
		assert.NotNil(t, v.Banana)

		if s, err := v.Apple.Action(); assert.NoError(t, err) {
			assert.Equal(t, "apple", s)
		}

		if s, err := v.Banana.Action(); assert.NoError(t, err) {
			assert.Equal(t, "banana", s)
		}
	})

	t.Run("nil value", func(t *testing.T) {
		assert.Panics(t, func() { InjectAllClient(nil) })
	})

	t.Run("not pointer", func(t *testing.T) {
		var v struct{}
		assert.Panics(t, func() { InjectAllClient(v) })
	})

	t.Run("not struct by literal", func(t *testing.T) {
		var v int
		assert.Panics(t, func() { InjectAllClient(v) })
	})

	t.Run("not struct by pointer", func(t *testing.T) {
		var v int
		assert.Panics(t, func() { InjectAllClient(&v) })
	})

	t.Run("literal struct", func(t *testing.T) {
		var v struct {
			Apple    AppleClient `grpc_client_inject:"apple"`
			Embedded struct {
				Banana BananaClient `grpc_client_inject:"banana"`
			}
		}

		assert.NotPanics(t, func() { InjectAllClient(&v) })
		assert.NotNil(t, v.Apple)
		assert.NotNil(t, v.Embedded.Banana)
	})

	t.Run("literal pointer", func(t *testing.T) {
		var v struct {
			Apple    AppleClient `grpc_client_inject:"apple"`
			Embedded *struct {
				Banana BananaClient `grpc_client_inject:"banana"`
			}
		}

		assert.NotPanics(t, func() { InjectAllClient(&v) })
		assert.NotNil(t, v.Apple)
		assert.NotNil(t, v.Embedded)
		assert.NotNil(t, v.Embedded.Banana)
	})

	t.Run("embedded struct", func(t *testing.T) {
		type Embedded struct {
			Banana BananaClient `grpc_client_inject:"banana"`
		}
		type Fruits struct {
			Embedded
			Apple AppleClient `grpc_client_inject:"apple"`
		}

		var v Fruits

		assert.NotPanics(t, func() { InjectAllClient(&v) })
		assert.NotNil(t, v.Apple)
		assert.NotNil(t, v.Embedded.Banana)
	})

	t.Run("embedded pointer", func(t *testing.T) {
		type Embedded struct {
			Banana BananaClient `grpc_client_inject:"banana"`
		}
		type Fruits struct {
			*Embedded
			Apple AppleClient `grpc_client_inject:"apple"`
		}

		var v Fruits

		assert.NotPanics(t, func() { InjectAllClient(&v) })
		assert.NotNil(t, v.Apple)
		assert.NotNil(t, v.Embedded)
		assert.NotNil(t, v.Embedded.Banana)
	})

	t.Run("bad name", func(t *testing.T) {
		var v struct {
			Apple  AppleClient  `grpc_client_inject:"apple"`
			Banana BananaClient `grpc_client_inject:"banana"`
			Cherry BananaClient `grpc_client_inject:"cherry"`
		}

		assert.Panics(t, func() { InjectAllClient(&v) })
	})

	t.Run("not interface", func(t *testing.T) {
		var v struct {
			Apple  *DummyClient `grpc_client_inject:"apple"`
			Banana BananaClient `grpc_client_inject:"banana"`
		}

		assert.NotPanics(t, func() { InjectAllClient(&v) })
		assert.NotNil(t, v.Banana)
	})

	t.Run("unexported", func(t *testing.T) {
		var v struct {
			Apple  AppleClient  `grpc_client_inject:"apple"`
			banana BananaClient `grpc_client_inject:"banana"`
		}

		assert.NotPanics(t, func() { InjectAllClient(&v) })
		assert.NotNil(t, v.Apple)
		assert.Nil(t, v.banana)
	})

	t.Run("untag", func(t *testing.T) {
		var v struct {
			Apple  AppleClient
			Banana BananaClient `grpc_client_inject:"banana"`
		}

		assert.NotPanics(t, func() { InjectAllClient(&v) })
		assert.Nil(t, v.Apple)
		assert.NotNil(t, v.Banana)
	})

	t.Run("bad tag", func(t *testing.T) {
		var v struct {
			Apple  AppleClient  `grpc_cilent_inject:"apple"`
			Banana BananaClient `grpc_client_inject:"banana"`
		}

		assert.NotPanics(t, func() { InjectAllClient(&v) })
		assert.Nil(t, v.Apple)
		assert.NotNil(t, v.Banana)
	})
}
