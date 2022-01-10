package handler_test

import (
	"github.com/zhoubiao2019/monitor-gateway/conf"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"

	"github.com/zhoubiao2019/monitor-gateway/model"
	"github.com/zhoubiao2019/monitor-gateway/router"
)

func initEnv(t *testing.T) *assert.Assertions {
	conf.InitConfig("../conf/config.yaml")
	//model.InitDB()

	return assert.New(t)
}

func TestRegisterJob(t *testing.T) {
	initEnv(t)
	guard := monkey.Patch(model.CreateJob, func(_ *model.Job) error {
		return nil
	})
	defer guard.Unpatch()

	engine := gin.Default()
	router.Init(engine)

	param := &model.Job{
		Name:             "job_test1",
		Status:           "healthy",
		NextScheduleTime: time.Now().Add(10 * time.Hour).Unix(),
		LastScheduleTime: time.Now().Unix(),
		Description:      "this is a test job",
	}
	paramRaw, _ := jsoniter.MarshalToString(param)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/v1/job/register", strings.NewReader(paramRaw))
	if err != nil {
		panic(err)
	}
	engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// todo: mock some error cases
}

func TestStatusReport(t *testing.T) {
	engine := gin.Default()
	router.Init(engine)

	guard := monkey.Patch(model.CreateJobLog, func(_ *model.JobLog) error {
		return nil
	})
	defer guard.Unpatch()

	param := &model.JobLog{
		JobName:          "job_test1",
		Result:           model.JobResultSuccess,
		NextScheduleTime: time.Now().Add(10 * time.Hour).Unix(),
		ScheduleTime:     time.Now().Unix(),
	}
	paramRaw, _ := jsoniter.MarshalToString(param)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/job/status_report", strings.NewReader(paramRaw))
	engine.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// todo: mock some error cases
}

func TestListJob(t *testing.T) {
	engine := gin.Default()
	router.Init(engine)

	guard := monkey.Patch(model.GetJobList, func() ([]*model.Job, error) {
		list := []*model.Job{
			&model.Job{
				ID:               1,
				Name:             "job_name_test002",
				Status:           model.JobStatusHealthy,
				NextScheduleTime: time.Now().Add(10 * time.Hour).Unix(),
				LastScheduleTime: time.Now().Unix(),
				Description:      "this is a test",
			},
			&model.Job{
				ID:               2,
				Name:             "job_name_test002",
				Status:           model.JobStatusUnhealthy,
				NextScheduleTime: time.Now().Add(10 * time.Hour).Unix(),
				LastScheduleTime: time.Now().Unix(),
				Description:      "this is a test",
			},
		}
		return list, nil
	})
	defer guard.Unpatch()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/job/list", nil)
	engine.ServeHTTP(w, req)

	//body, _ := ioutil.ReadAll(w.Body)
	//var list []*model.Job
	//jsoniter.Unmarshal(body, &list)

	assert.Equal(t, 200, w.Code)
	//todo: check body
}
