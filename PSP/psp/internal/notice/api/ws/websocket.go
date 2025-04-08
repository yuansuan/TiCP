package ws

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/service/mq"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

const (
	DefaultConnSize = 100
)

type webSocketService struct {
	mutex       sync.Mutex
	consumer    *mq.KafkaConsumer
	upgrade     *websocket.Upgrader
	connections []*websocket.Conn
}

// NewWebSocketService 创建WebSocket服务
func NewWebSocketService() (*webSocketService, error) {
	consumer, err := mq.NewKafkaConsumer(common.NoticeWebsocketTopic, common.NoticeWebsocketGroup)
	if err != nil {
		return nil, err
	}

	return &webSocketService{
		consumer: consumer,
		upgrade: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		connections: make([]*websocket.Conn, 0, DefaultConnSize),
	}, nil
}

// WebsocketNotice websocket消息通知
func (s *webSocketService) WebsocketNotice(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	conn, err := s.upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Errorf("upgrade websocket connection err: %v", err)
		ginutil.Error(ctx, errcode.ErrInternalServer, errcode.MsgInternalServer)
		return
	}

	s.mutex.Lock()
	s.connections = append(s.connections, conn)
	s.mutex.Unlock()

	go s.pingHandle(ctx, conn)

	for {
		message, err := s.consumer.ReadByteMessage(context.Background())
		if err != nil {
			logger.Errorf("read kafka message err: %v", err)
			ginutil.Error(ctx, errcode.ErrInternalServer, errcode.MsgInternalServer)
			return
		}

		s.broadcastMessage(ctx, message)
	}
}

func (s *webSocketService) removeConnection(ctx context.Context, conn *websocket.Conn) {
	defer conn.Close()

	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i, c := range s.connections {
		if c == conn {
			s.connections = append(s.connections[:i], s.connections[i+1:]...)
			break
		}
	}
}

func (s *webSocketService) broadcastMessage(ctx context.Context, message []byte) {
	logger := logging.GetLogger(ctx)
	logger.Infof("the message broadcast to [%v] websocket connections", len(s.connections))

	deadConnections := make([]*websocket.Conn, 0, len(s.connections))
	for _, conn := range s.connections {
		if err := s.writeMessage(conn, message); err != nil {
			logger.Infof("remove the [%v] websocket connection that cannot write message, reason: %v", conn.RemoteAddr().String(), err)
			deadConnections = append(deadConnections, conn)
		} else {
			logger.Debugf("[%v] websocket connection write message success", conn.RemoteAddr().String())
		}
	}

	for _, connection := range deadConnections {
		s.removeConnection(ctx, connection)
	}
}

func (s *webSocketService) pingHandle(ctx context.Context, conn *websocket.Conn) {
	logger := logging.GetLogger(ctx)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			logger.Infof("websocket reade ping message err: %v", err)
			s.removeConnection(ctx, conn)
			return
		}

		if err := s.writeMessage(conn, []byte("Pong")); err != nil {
			logger.Errorf("websocket write pong message err: %v", err)
			return
		}
	}
}

func (s *webSocketService) writeMessage(conn *websocket.Conn, message []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return conn.WriteMessage(websocket.TextMessage, message)
}
