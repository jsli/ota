package ota_job

import (
	"fmt"
	"github.com/jsli/gtbox/file"
	"github.com/jsli/gtbox/pathutil"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"github.com/jsli/ota/radio/app/policy"
	"github.com/robfig/revel"
	"time"
)

type ReleaseCreationJob struct {
}

func (rcj *ReleaseCreationJob) Run() {
	tag := "JOB_D: "
	revel.INFO.Println(tag, "running")
	dal, err := models.NewDal()
	if err != nil {
		revel.ERROR.Println(tag, err)
		return
	}
	defer dal.Close()

	task, err := models.PopOneCreationTask(dal)
	if err != nil {
		revel.ERROR.Println(tag, err)
		return
	}

	if task != nil {
		task.Flag = ota_constant.FLAG_CREATING
		task.ModifiedTs = time.Now().Unix()
		task.Update(dal)

		root_path := fmt.Sprintf("%s%s/", ota_constant.TMP_FILE_ROOT, policy.GenerateRandFileName())
		pathutil.MkDir(root_path)
		defer file.DeleteDir(root_path)
		revel.INFO.Println(tag, "Processing task : ", task.UpdateRequest)
		revel.INFO.Println(tag, "TEMP dir : ", root_path)

		release, err := policy.GenerateOtaPackage(dal, task, root_path)
		if err != nil {
			revel.ERROR.Println(tag, "Failed: ", err)
			revel.INFO.Println(tag, "Failed, task id= ", task.Id, " error msg: ", err)
			if task.RetryCount >= 5 {
				task.Flag = ota_constant.FLAG_DROPPED
			} else {
				task.RetryCount = task.RetryCount + 1
				task.Flag = ota_constant.FLAG_INIT
				revel.INFO.Println(tag, "Failed, task id= ", task.Id, " retry: ", task.RetryCount)
			}
			task.ModifiedTs = time.Now().Unix()
			_, uerr = task.Update(dal)
			revel.ERROR.Println(tag, " task UPDATE Failed: ", uerr)
			return
		}

		task.ReleaseId = release.Id
		task.Flag = ota_constant.FLAG_CREATED
		task.ModifiedTs = time.Now().Unix()
		_, uerr := task.Update(dal)
		if uerr != nil {
			revel.ERROR.Println(tag, task, " task UPDATE Failed: ", uerr)
		}
	}
}
