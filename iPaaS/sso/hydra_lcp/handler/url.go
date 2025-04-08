package handler

import (
	"os"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getLoginURL(c *gin.Context, challenge string) string {
	// get frontend server from env
	fs := os.Getenv("FRONTEND_SERVER_ENDPOINT")
	if "" == fs {
		fs = h.conf.FrontendServer
	}
	return fs + "/#/login?login_challenge=" + challenge
}

func (h *Handler) getLogoutURL(c *gin.Context, challenge string) string {
	// get frontend server from env
	fs := os.Getenv("FRONTEND_SERVER_ENDPOINT")
	if "" == fs {
		fs = h.conf.FrontendServer
	}
	return fs + "/#/logout?logout_challenge=" + challenge
}

func (h *Handler) getLoginURLPortal(c *gin.Context, challenge string) string {
	// get frontend server from env
	fs := os.Getenv("FRONTEND_PORTAL_SERVER_ENDPOINT")
	if "" == fs {
		fs = h.conf.FrontendPortalServer
	}
	return fs + "/#/login?login_challenge=" + challenge
}

func (h *Handler) getLogoutURLPortal(c *gin.Context, challenge string) string {
	// get frontend server from env
	fs := os.Getenv("FRONTEND_PORTAL_SERVER_ENDPOINT")
	if "" == fs {
		fs = h.conf.FrontendPortalServer
	}
	return fs + "/#/logout?logout_challenge=" + challenge
}
