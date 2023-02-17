package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// monitor "github.com/ljahier/uptime-monitor/pkg/monitor"
)

func RunWebServer() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8081")
}
