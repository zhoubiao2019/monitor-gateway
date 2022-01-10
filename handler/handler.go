package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhoubiao2019/monitor-gateway/db"
	"github.com/zhoubiao2019/monitor-gateway/util/log"
	"time"

	"github.com/zhoubiao2019/monitor-gateway/model"
)

func RegisterJob(ctx *gin.Context) {
	job := new(model.Job)
	if err := ctx.BindJSON(job); err != nil {
		ctx.String(500, "bind json failed,"+err.Error())
		return
	}

	session := db.DBer.NewSession()
	defer session.Close()

	if err := model.CreateJob(session, job); err != nil {
		ctx.String(500, "create_job failed,"+err.Error())
		return
	}
	ctx.String(200, "ok")
}

func StatusReport(ctx *gin.Context) {
	jobLog := new(model.JobLog)
	if err := ctx.BindJSON(jobLog); err != nil {
		ctx.String(500, "bind json failed,"+err.Error())
		return
	}

	session := db.DBer.NewSession()
	defer session.Close()

	jobLog.ClientIP = ctx.ClientIP()

	if err := model.CreateJobLog(session, jobLog); err != nil {
		ctx.String(500, "create_job_log failed,"+err.Error())
		return
	}

	ctx.String(200, "ok")
}

func ListJob(ctx *gin.Context) {
	session := db.DBer.NewSession()
	defer session.Close()

	list, err := model.GetJobList(session)
	if err != nil {
		ctx.String(500, "get_job_list failed")
		return
	}

	ctx.JSON(200, list)
}

func RefreshJobStatus() {
	session := db.DBer.NewSession()
	defer session.Close()

	list, err := model.GetJobList(session)
	if err != nil {
		log.Errorw("refresh failed", "err", err.Error())
		return
	}

	for _, job := range list {
		if time.Now().Unix() < job.NextScheduleTime {
			continue
		}

		job := &model.Job{ID: job.ID, Status: model.JobStatusUnhealthy}
		if err := model.UpdateJob(session, job); err != nil {
			log.Errorw("update job failed", "err", err.Error())
		}
	}
}
