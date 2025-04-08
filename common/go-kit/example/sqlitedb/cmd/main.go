package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/example/sqlitedb/internal/dao"
	_ "modernc.org/sqlite"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	_http "github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
)

func init() {
	_ = os.Setenv("YS_MODE", "local")
	_ = os.Setenv("YS_LOG_LEVEL", "DEBUG")

	boot.MW.UseDefaultDatabase(middleware.DatabaseSQLite)
}

func main() {
	server := boot.Default()

	server.Register(func(*_http.Driver) {
		sess := boot.MW.DefaultSession(context.Background())
		defer func() { _ = sess.Close() }()

		result, err := sess.Query(`create table if not exists key_map (key text, value text);`)
		if err != nil {
			panic(err)
		}
		log.Println("Create Table Result: ", result)
	})

	server.Register(func(server *_http.Driver) {
		server.Any("/map/put", func(ctx *gin.Context) {
			if key, ok := ctx.GetQuery("key"); ok {
				if val, ok := ctx.GetQuery("val"); ok {
					if err := dao.Set(ctx, key, val); err == nil {
						ctx.JSON(http.StatusOK, gin.H{
							"key": key,
							"val": val,
						})
					} else {
						ctx.JSON(http.StatusBadRequest, gin.H{
							"error": err.Error(),
						})
					}

					return
				}
			}

			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid user input",
			})
		})

		server.Any("/map/get", func(ctx *gin.Context) {
			if key, ok := ctx.GetQuery("key"); ok {
				if m, err := dao.Get(ctx, key); err == nil {
					ctx.JSON(http.StatusOK, gin.H{
						"key": m.Key,
						"val": m.Value,
					})
				} else {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
				}

				return
			}

			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid user input",
			})
		})
	}).Run()
}
