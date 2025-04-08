package rdpgo

type BaseRequest struct {
	PrivateIP string `query:"PrivateIP"`
	RequestID string `header:"x-ys-request-id"`
}
