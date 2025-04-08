package iamserver

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
)

type Server struct {
}

func NewApiServer() *Server {
	return &Server{}
}

func (s *Server) ToJob(c *gin.Context) {
	s.To(c, config.GetConfig().ApiServerT.Job)
}

func (s *Server) ToCloudApp(c *gin.Context) {
	s.To(c, config.GetConfig().ApiServerT.CloudApp)
}

func (s *Server) ToAccBill(c *gin.Context) {
	s.To(c, config.GetConfig().ApiServerT.AccountBill)
}
func (s *Server) ToLicManage(c *gin.Context) {
	s.To(c, config.GetConfig().ApiServerT.LicManager)
}

func (s *Server) ToMerchandise(c *gin.Context) {
	s.To(c, config.GetConfig().ApiServerT.Merchandise)
}

func (s *Server) To(c *gin.Context, targetURL string) {
	targetURLParsed, err := url.Parse(targetURL)
	if err != nil {
		logging.Default().Warnf("TranFail, Error: %s, Url: %s", err.Error(), targetURL)
		common.InternalServerError(c, "")
		return
	}
	originalURL := c.Request.URL
	originalURL.Scheme = targetURLParsed.Scheme
	originalURL.Host = targetURLParsed.Host
	req, err := http.NewRequest(c.Request.Method, originalURL.String(), c.Request.Body)
	if err != nil {
		logging.Default().Warnf("TranFail, Error: %s, Url: %s", err.Error(), targetURL)
		common.InternalServerError(c, "")
		return
	}
	req.Header = make(http.Header)
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	// context deadline exceeded (Client.Timeout exceeded while awaiting headers)
	client := &http.Client{Timeout: time.Second * 30} // avoid: Client.Timeout
	resp, err := client.Do(req)
	if err != nil {
		logging.Default().Warnf("TranFail, Error: %s, Url: %s", err.Error(), targetURL)
		common.ErrorRespWithAbort(c, http.StatusBadGateway, "InternalServerError", "")
		return
	}
	defer closeResponse(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Default().Warnf("TranFail, Error: %s, Url: %s", err.Error(), targetURL)
		common.InternalServerError(c, "")
		return
	}
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}
	c.Header(common.RequestIDKey, common.GetRequestID(c))
	c.Status(resp.StatusCode)
	c.Writer.Write(body)
}

// closeResponse close non nil response with any response Body.
// convenient wrapper to drain any remaining data on response body.
// Subsequently this allows golang http RoundTripper to re-use
// - the same connection for future requests.
func closeResponse(resp *http.Response) {
	// Callers should close resp.Body when done reading from it.
	// If resp.Body is not closed, the Client's underlying RoundTripper
	// (typically Transport) may not be able to re-use a persistent TCP
	// connection to the server for a subsequent "keep-alive" request.
	if resp != nil && resp.Body != nil {
		// Drain any remaining Body and then close the connection.
		// Without this closing connection would disallow re-using
		// the same connection for future uses.
		//  - http://stackoverflow.com/a/17961593/4465767
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
