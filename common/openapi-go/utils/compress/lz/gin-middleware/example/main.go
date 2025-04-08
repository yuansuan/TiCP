package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	middleware "github.com/yuansuan/ticp/common/openapi-go/utils/compress/lz/gin-middleware"
)

/*
curl -v http://localhost:8080/ping \
-H "Accept-Encoding: lz"  \
--output -
*/

// * Uses proxy env variable http_proxy == 'http://127.0.0.1:7890'
// *   Trying 127.0.0.1:7890...
// * Connected to 127.0.0.1 (127.0.0.1) port 7890 (#0)
// > GET http://localhost:8080/ping HTTP/1.1
// > Host: localhost:8080
// > User-Agent: curl/7.84.0
// > Accept: */*
// > Proxy-Connection: Keep-Alive
// > Accept-Encoding: lz
// >
// * Mark bundle as not supporting multiuse
// < HTTP/1.1 200 OK
// < Content-Length: 280
// < Connection: keep-alive
// < Content-Encoding: lz
// < Content-Type: text/plain; charset=utf-8
// < Date: Tue, 01 Nov 2022 07:48:49 GMT
// < Keep-Alive: timeout=4
// < Proxy-Connection: keep-alive
// < Vary: Accept-Encoding
// <
// * Connection #0 to host 127.0.0.1 left intact
// {garbled/binary output}

func main() {
	r := gin.Default()
	r.Use(middleware.Lz())
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})

	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
