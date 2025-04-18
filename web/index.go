package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nahidhasan98/remind-name/config"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.gohtml", gin.H{
		"Title":         "Remind The Names of Allah (SWT)",
		"AppName":       "Remind Name",
		"VersionString": config.VERSION_STRING,
	})
}
