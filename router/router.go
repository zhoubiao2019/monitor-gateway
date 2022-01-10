package router

import (
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zhoubiao2019/monitor-gateway/handler"
	"github.com/zhoubiao2019/monitor-gateway/util/log"
)

func Init(g *gin.Engine) {
	g.Use(Log)

	g.POST("/api/v1/job/register", handler.RegisterJob)
	g.POST("/api/v1/job/status_report", handler.StatusReport)
	g.GET("/api/v1/job/list", handler.ListJob)
}

func Log(c *gin.Context) {
	st := time.Now()
	c.Next()
	latency := time.Since(st).Seconds() * 1000

	var (
		errCode = "0"
	)
	if code, exist := c.Get("err_code"); exist {
		errCode = strconv.Itoa(int(code.(int32)))
	}
	req, err := httputil.DumpRequest(c.Request, true)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Infow("deploy gateway access log",
		"user", c.GetString("sessionValue"),
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
		"request", string(req),
		"status", c.Writer.Status(),
		"errCode", errCode,
		"latency(ms)", latency,
	)
}
