package model

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
)

type JobStatus string
type JobExecuteResult string

const (
	JobStatusHealthy   JobStatus = "healthy"
	JobStatusUnhealthy JobStatus = "unhealthy"

	JobResultSuccess JobExecuteResult = "success"
	JobResultFailed  JobExecuteResult = "failed"
)

type Job struct {
	ID               int64 `xorm:"pk autoincr"`
	Name             string
	Status           JobStatus
	NextScheduleTime int64
	LastScheduleTime int64
	LastScheduleID   int64
	Description      string

	UpdatedAt time.Time `xorm:"updated"`
	CreatedAt time.Time `xorm:"created"`
}

type JobLog struct {
	ID               int64 `xorm:"pk autoincr"`
	JobID            int64
	JobName          string
	Result           JobExecuteResult
	Detail           string
	ClientIP         string
	ScheduleTime     int64
	NextScheduleTime int64

	CreatedAt time.Time `xorm:"created"`
}

func CreateJob(session *xorm.Session, job *Job) error {
	_, err := session.InsertOne(job)
	return err
}

func UpdateJob(session *xorm.Session, job *Job) error {
	_, err := session.Table(&Job{}).ID(job.ID).Update(job)
	return err
}

func GetJobList(session *xorm.Session) ([]*Job, error) {
	var list []*Job
	err := session.Table(&Job{}).Find(&list)
	return list, err
}

func GetJobByName(session *xorm.Session, name string) (*Job, error) {
	job := new(Job)
	isExist, err := session.Table(&Job{}).Where("name = ?", name).Get(job)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, fmt.Errorf("job not found")
	}

	return job, nil
}

func CreateJobLog(session *xorm.Session, jobLog *JobLog) error {
	job, err := GetJobByName(session, jobLog.JobName)
	if err != nil {
		return err
	}

	jobLog.JobID = job.ID
	if _, err = session.InsertOne(jobLog); err != nil {
		return err
	}

	param := new(Job)
	if jobLog.Result == JobResultSuccess {
		param.Status = JobStatusHealthy
	} else {
		param.Status = JobStatusUnhealthy
	}
	param.LastScheduleTime = jobLog.ScheduleTime
	param.NextScheduleTime = jobLog.NextScheduleTime
	if _, err := session.Table(&Job{}).ID(job.ID).Update(param); err != nil {
		return err
	}

	return nil
}
