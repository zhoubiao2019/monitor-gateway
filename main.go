package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zhoubiao2019/monitor-gateway/conf"
	"github.com/zhoubiao2019/monitor-gateway/handler"
	"github.com/zhoubiao2019/monitor-gateway/model"
	"github.com/zhoubiao2019/monitor-gateway/router"
	"github.com/zhoubiao2019/monitor-gateway/util/grace"
	"github.com/zhoubiao2019/monitor-gateway/util/log"
)

func main() {
	confPath := flag.String("f", "./conf/config.yaml", "config file path")
	flag.Parse()
	conf.InitConfig(*confPath)
	initLogger()
	model.InitDB()

	gin.SetMode(conf.GlobalConfig.Gin.Mode)
	engine := gin.Default()
	router.Init(engine)
	srv := &http.Server{
		Addr:    ":" + conf.GlobalConfig.Gin.Port,
		Handler: engine,
	}

	var grace grace.Grace
	grace.Register(func() {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		if err := srv.Shutdown(ctx); err != nil {
			log.Errorf("Graceful shutting down failed. err=%v", err)
		}
	})
	grace.Run(func() error {
		go func() { StartRefreshJob() }()
		return srv.ListenAndServe()
	})
}

func StartRefreshJob() {
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			handler.RefreshJobStatus()
		}
	}
}

func initLogger() {
	logConf := conf.GlobalConfig.Logger
	logger := log.New(
		log.WithLevel(logConf.Level),
		log.WithFileLog(logConf.FileLog),
		log.WithFilePath(logConf.Path),
		log.WithConsoleLog(logConf.ConsoleLog),
	)
	log.SetDefault(logger)
}
