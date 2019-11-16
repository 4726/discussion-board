package common

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func AddMonitorHandler(engine *gin.Engine) {
	engine.Any("/metrics", gin.WrapH(promhttp.Handler()))
}