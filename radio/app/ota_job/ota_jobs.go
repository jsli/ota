package ota_job

import (
	"fmt"
	"github.com/jsli/gtbox/file"
	"github.com/jsli/gtbox/pathutil"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"github.com/jsli/ota/radio/app/policy"
	"github.com/robfig/revel"
)

type ReleaseCreationJob struct {
}

func (rcj *ReleaseCreationJob) Run() {
	revel.INFO.Println("Create job : running")
	dal, err := models.NewDal()
	if err != nil {
		revel.ERROR.Println("Create job error: ", err)
		return
	}
	defer dal.Close()

	task, err := models.PopOneCreationTask(dal)
	if err != nil {
		revel.ERROR.Println("Create job error: ", err)
		return
	}

	if task != nil {
		root_path := fmt.Sprintf("%s%s/", ota_constant.TMP_FILE_ROOT, policy.GenerateRandFileName())
		pathutil.MkDir(root_path)
		defer file.DeleteDir(root_path)
		revel.INFO.Println("Create job : processing task : ", task, " ----- tmp dir : ", root_path)

		task.Flag = ota_constant.FLAG_CREATING
		task.Update(dal)

		release, err := policy.GenerateOtaPackage(dal, task, root_path)
		if err != nil {
			task.Flag = ota_constant.FLAG_CREATE_FAILED
			task.Update(dal)
			revel.INFO.Println("Create job : processing task failed: ", task)
		}

		task.ReleaseId = release.Id
		task.Flag = ota_constant.FLAG_CREATED
		task.Update(dal)
	}
}