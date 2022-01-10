package model

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"

	"github.com/zhoubiao2019/monitor-gateway/conf"
	"github.com/zhoubiao2019/monitor-gateway/util/log"
)

var DBer *xorm.Engine

func InitDB() {
	var err error
	DBer, err = xorm.NewEngine("mysql", conf.GlobalConfig.DB.DSN)
	if err != nil {
		log.Fatal(err.Error())
	}

	DBer.SetMaxIdleConns(conf.GlobalConfig.DB.MaxIdle)
	DBer.SetTableMapper(core.GonicMapper{})
	DBer.SetColumnMapper(core.GonicMapper{})
	DBer.DatabaseTZ = time.UTC
	DBer.TZLocation = time.UTC
	DBer.ShowExecTime(true)
	DBer.ShowSQL(true)

	err = DBer.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = DBer.Sync2(new(Job), new(JobLog))
	if err != nil {
		log.Fatalf("db sync failed:%v", err)
	}
	log.Info("db connected")
}
