package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zhoubiao2019/monitor-gateway/model"
	"github.com/zhoubiao2019/monitor-gateway/util/log"
)

func RegisterJob(ctx *gin.Context) {
	job := new(model.Job)
	if err := ctx.BindJSON(job); err != nil {
		ctx.String(500, "bind json failed,"+err.Error())
		return
	}

	if err := model.CreateJob(job); err != nil {
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

	jobLog.ClientIP = ctx.ClientIP()
	if err := model.CreateJobLog(jobLog); err != nil {
		ctx.String(500, "create_job_log failed,"+err.Error())
		return
	}

	ctx.String(200, "ok")
}

func ListJob(ctx *gin.Context) {
	list, err := model.GetJobList()
	if err != nil {
		ctx.String(500, "get_job_list failed")
		return
	}

	ctx.JSON(200, list)
}

func RefreshJobStatus() {
	list, err := model.GetJobList()
	if err != nil {
		log.Errorw("refresh failed", "err", err.Error())
		return
	}

	for _, job := range list {
		if time.Now().Unix() < job.NextScheduleTime {
			continue
		}

		job := &model.Job{ID: job.ID, Status: model.JobStatusUnhealthy}
		if err := model.UpdateJob(job); err != nil {
			log.Errorw("update job failed", "err", err.Error())
		}
	}
}
