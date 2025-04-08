package gin

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/rdpgo/guacamole"
	"github.com/yuansuan/ticp/rdpgo/jwt"
)

func encodeJWT(c *gin.Context) {
	req := new(guacamole.ConnectArgsInToken)
	if err := c.Bind(req); err != nil {
		c.Data(400, "text/plain; charset=utf-8", []byte("failed parse request"))
		return
	}

	bs, _ := json.Marshal(req)

	fmt.Println(req)

	jwtStr, err := jwt.Encode(string(bs))
	if err != nil {
		c.Data(500, "text/plain; charset=utf-8", []byte(err.Error()))
		return
	}

	c.Data(200, "text/plain; charset=utf-8", []byte(jwtStr))
}

//go:embed web/*
var fs embed.FS

func Setup(e *gin.Engine) {
	e.GET("/jwt", func(c *gin.Context) {
		dirs, _ := fs.ReadDir("web")
		for _, v := range dirs {
			fmt.Println(v.Name())
		}

		data, err := fs.ReadFile("web/jwt.html")
		if err != nil {
			c.Data(500, "text/plain; charset=utf-8", []byte(err.Error()))
			return
		}

		c.Data(http.StatusOK, "text/html;", data)
	})

	e.POST("/jwt", encodeJWT)
}
