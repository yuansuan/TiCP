package api

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/notice/api/ws"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/service/impl"
)

type apiRoute struct {
	emailService   service.EmailService
	messageService service.MessageService
}

// NewAPIRoute 创建api服务
func NewAPIRoute() (*apiRoute, error) {
	emailService, err := impl.NewEmailService()
	if err != nil {
		return nil, err
	}
	messageService, err := impl.NewMessageService()
	if err != nil {
		return nil, err
	}

	return &apiRoute{
		emailService:   emailService,
		messageService: messageService,
	}, nil
}

// InitAPI 初始化API服务
func InitAPI(drv *http.Driver) {
	logger := logging.Default()

	api, err := NewAPIRoute()
	if err != nil {
		logger.Errorf("new api service err: %v", err)
		panic(err)
	}

	ws, err := ws.NewWebSocketService()
	if err != nil {
		logger.Errorf("new websocket service err: %v", err)
		panic(err)
	}

	group := drv.Group("/api/v1")
	{
		noticeGroup := group.Group("/notice")

		// message
		noticeGroup.PUT("/read", api.ReadMessage)
		noticeGroup.PUT("/readAll", api.ReadAllMessage)
		noticeGroup.GET("/count", api.MessageCount)
		noticeGroup.POST("/list", api.MessageList)
		noticeGroup.POST("/producer", api.SendWebsocketMessage)

		// email
		noticeGroup.GET("/email", api.GetEmail)
		noticeGroup.POST("/email", api.SetEmail)
		noticeGroup.POST("/email/send", api.SendEmail)
		noticeGroup.POST("/email/testSend", api.TestSendEmail)
	}

	wsGroup := drv.Group("/ws/v1")
	{
		noticeGroup := wsGroup.Group("/notice")
		// websocket api
		noticeGroup.GET("/consumer", ws.WebsocketNotice)
	}
}
