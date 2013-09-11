package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/log"
	"github.com/jsli/ota/radio/app/models"
	"github.com/jsli/ota/radio/app/utils"
	"github.com/robfig/revel"
	"os"
	"time"
)

type Radio struct {
	*revel.Controller
}

func (c Radio) Index() revel.Result {
	return c.Render()
}

func (c Radio) Download() revel.Result {
	tag := "Download"
	var model, t, bldv, cpv string
	c.Params.Bind(&model, "model")
	c.Params.Bind(&t, "type")
	c.Params.Bind(&bldv, "bldv")
	c.Params.Bind(&cpv, "cpv")
	log.Log(tag, fmt.Sprintf("params : model=%s, type=%s, bldv=%s, cpv=%s\n", model, t, bldv, cpv))

	var row *sql.Row
	var query string
	//	cpv_int := utils.ConvertVersion2Int64(cpv)
	dal, err := models.NewDal(models.DRIVER, models.DNS)
	if err != nil {
		log.Log(tag, fmt.Sprintf("New dal error: %s\n", err))
		return c.Redirect(Radio.Index)
	}
	defer dal.Close()

	var img_path string
	if t == "single" {
		query = fmt.Sprintf("SELECT id, release_name FROM radio_release where (model='%s' and build_version='%s' and single_version='%s')", model, bldv, cpv)
		img_path = fmt.Sprintf("%s%s/single/%s/radio.img", constant.UPLOAD_ROOT_DIR, model, cpv)
	} else if t == "dsds" {
		query = fmt.Sprintf("SELECT id, release_name FROM radio_release where (model='%s' and build_version='%s' and dsds_version='%s')", model, bldv, cpv)
		img_path = fmt.Sprintf("%s%s/dsds/%s/radio.img", constant.UPLOAD_ROOT_DIR, model, cpv)
	} else {
		log.Log(tag, fmt.Sprintf("wrong params type=%s, should be 'single' or 'dsds' only!\n", model, t, bldv, cpv))
		return c.Redirect(Radio.Index)
	}
	log.Log(tag, fmt.Sprintf("query : %s \n", query))

	row = dal.Link.QueryRow(query)
	var id int = -1
	err = row.Scan(&id)

	//we didn't generate update package for this version
	if err != nil {
		log.Log(tag, fmt.Sprintf("query rows error: %s\n", err))
		log.Log(tag, fmt.Sprintf("generate ota package to : %s%s\n", err))
		product, err := models.FindProduct(dal, model)
		if err != nil {
			log.Log(tag, fmt.Sprintf("generate ota package abort, illegal model=%s, errro: \n", model, err))
			return c.Redirect(Radio.Index)
		}
		fmt.Println(product)

		root := fmt.Sprintf("%stmp_%d/", constant.TEMP_DIR, time.Now().Unix())
		utils.TouchDir(root, 0)
		full_pkg_path := fmt.Sprintf("%s%s/%s-target_files-%s%s", constant.FULL_BUILD_DIR, bldv, model, bldv, constant.ZIP_SUFFIX)
		unzip_dir := root + "unzip/"
		utils.TouchDir(unzip_dir, 0)
		utils.CopyFileWithPath(full_pkg_path, unzip_dir+"full_pkg.zip")
		utils.ExtractZipFile(unzip_dir+"full_pkg.zip", unzip_dir) // (1) unzip file

		tmp_dir := root + "tmp/"
		utils.TouchDir(tmp_dir, 0)
		//(2)copy related files
		for _, f := range constant.COPY_FILE_LIST {
			utils.CopyFileWithPath(unzip_dir+f, tmp_dir+f)
		}
		for _, d := range constant.COPY_DIR_LIST {
			utils.CopyDirWithPath(unzip_dir+d, tmp_dir)
		}

		//(3) copy radio.img
		utils.CopyFileWithPath(img_path, tmp_dir+"radio.img")

		//(4)generate update_pkg.zip
		//archive all files, generate update package
		final_dir := root + "final/"
		update_pkg_dir := final_dir + "update_pkg/"
		utils.TouchDir(update_pkg_dir, 0)
		update_pkg_path := update_pkg_dir + "update_pkg.zip"
		utils.MakeZipFile(tmp_dir, update_pkg_path)

		ota_pkg_dir := final_dir + "ota_pkg/"
		ota_pkg_path := ota_pkg_dir + "update.zip"
		utils.TouchDir(ota_pkg_dir, 0)
		cmd_params := make([]string, 5)
		cmd_params[0] = constant.OTA_CMD_PARAM_PLATFORM_PREFIX + product.Platform //params[otasdk.PLATFORM].(string)
		cmd_params[1] = constant.OTA_CMD_PARAM_PRODUCT_PREFIX + product.Model     //params[otasdk.PRODUCT].(string)
		cmd_params[2] = constant.OTA_CMD_PARAM_OEM_PREFIX + product.Vendor        //params[otasdk.OEM].(string)
		cmd_params[3] = constant.OTA_CMD_PARAM_OUTPUT_PREFIX + ota_pkg_path
		cmd_params[4] = constant.OTA_CMD_PARAM_INPUT_PREFIX + update_pkg_path
		utils.GenerateOtaPackage("/home/manson/server/ota/new/radio/updatetool/updatemk", cmd_params)

		log.Log(tag, fmt.Sprintf("extract zip file %s to %s", full_pkg_path, unzip_dir))
	}

	log.Log(tag, fmt.Sprintf("query row result : id=%d\n", id))

	//	return c.Redirect(dir + "update.zip")
	return c.Redirect(Radio.Index)
}

func ParseParams(params *revel.Params) (map[string]string, error) {
	tag := "Radio.ParseParams"
	var model, t, bldv, cpv string
	params.Bind(&model, "model")
	if !utils.IsAvailableModel(model) {
		msg := fmt.Sprintf("illegal param model = %s\n", model)
		log.Log(tag, msg)
		return nil, errors.New(msg)
	}
	params.Bind(&t, "type")
	if !utils.IsAvailableType(t) {
		msg := fmt.Sprintf("illegal param type = %s\n", t)
		log.Log(tag, msg)
		return nil, errors.New(msg)
	}
	params.Bind(&cpv, "cpv")
	var cpv_extracted string
	var b bool
	if b, cpv_extracted = utils.IsAvailableCPV(cpv); !b {
		msg := fmt.Sprintf("illegal param cpv = %s\n", cpv)
		log.Log(tag, msg)
		return nil, errors.New(msg)
	}
	params.Bind(&bldv, "bldv")
	result := make(map[string]string)
	result["bldv"] = bldv
	result["model"] = model
	result["type"] = t
	result["cpv"] = cpv_extracted
	return result, nil
}

/*
 *URL Query parameters:
 *model: pxaxxx
 *type: single or dsds
 *bldv: build_version
 *cpv: cp version
 *
 */
func (c Radio) Query() revel.Result {
	tag := "Query"
	result := models.QueryResult{}
	extra := make(map[string]interface{})
	extra["api_version"] = "1.0"
	extra["error"] = ""
	result.Extra = extra
	params, err := ParseParams(c.Params)
	if err != nil {
		extra["error"] = fmt.Sprintf("%s", err)
		result.Data = nil
		return c.RenderJson(result)
	}
	log.Log(tag, fmt.Sprintf("params : %s\n", params))

	var rows *sql.Rows
	var query string
	cpv_int := utils.ConvertVersion2Int64(params["cpv"])
	dal, err := models.NewDal(models.DRIVER, models.DNS)
	if err != nil {
		log.Log(tag, fmt.Sprintf("New dal error: %s\n", err))
		extra["error"] = "server error!"
		return c.RenderJson(result)
	}
	defer dal.Close()

	if params["type"] == "single" {
		query = fmt.Sprintf("SELECT single_version_str FROM radio_image where model = '%s' and %s > %d", params["model"], "single_version", cpv_int)
	} else if params["type"] == "dsds" {
		query = fmt.Sprintf("SELECT dsds_version_str FROM radio_image where model = '%s' and %s > %d", params["model"], "dsds_version", cpv_int)
	} else {
		msg := fmt.Sprintf("wrong params type=%s, should be 'single' or 'dsds' only!\n", params["type"])
		log.Log(tag, msg)
		extra["error"] = msg
		return c.RenderJson(result)
	}
	log.Log(tag, fmt.Sprintf("query : %s \n", query))
	rows, err = dal.Link.Query(query)
	if err != nil {
		log.Log(tag, fmt.Sprintf("query rows error: %s\n", err))
		extra["error"] = "server error!"
		return c.RenderJson(result)
	}
	defer rows.Close()

	entries := make([]models.QueryResultEntry, 0, 10)
	for rows.Next() {
		var version string
		_ = rows.Scan(&version)
		entry := models.QueryResultEntry{params["model"], params["type"], version}
		entries = append(entries, entry)
	}
	result.Data = entries
	result.Count = len(entries)
	extra["error"] = nil
	log.Log(tag, fmt.Sprintf("query rows result: %s\n", result))
	return c.RenderJson(result)
}

func (c Radio) GetUpload() revel.Result {
	models := utils.MODEL_LIST
	return c.Render(models)
}

func (c Radio) PostUpload(upload *models.RadioImageFile) revel.Result {
	tag := "PostUpload"

	upload.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(Radio.PostUpload)
	}

	//convert version_str to version_int for using in db
	upload.Single.VersionInt = utils.ConvertVersion2Int64(upload.Single.VersionStr)
	upload.Dsds.VersionInt = utils.ConvertVersion2Int64(upload.Dsds.VersionStr)

	//mkdir for cp files
	//dir struct like this:
	//tree constant.UPLOAD_ROOT_DIR
	//├── pxa1t88ff_def
	//│   ├── dsds
	//│   │   └── 3.345.003 -> /home/manson/temp/test/CP/pxa1t88ff_def/single/2.456.003/
	//│   └── single
	//│       └── 2.456.003
	//│           ├── HL_TD_CP.bin
	//│           ├── HL_TD_DSDS_CP.bin
	//│           ├── HL_TD_M08_AI_A0_DSDS_Flash.bin
	//│           ├── HL_TD_M08_AI_A0_Flash.bin
	//│           └── radio.img
	//└── pxa986ff_def
	//    ├── dsds
	//    │   └── 3.345.004 -> /home/manson/temp/test/CP/pxa986ff_def/single/2.456.004/
	//    └── single
	//        └── 2.456.004
	//            ├── HL_TD_CP.bin
	//            ├── HL_TD_DSDS_CP.bin
	//            ├── HL_TD_M08_AI_A0_DSDS_Flash.bin
	//            ├── HL_TD_M08_AI_A0_Flash.bin
	//            └── radio.img

	single := fmt.Sprintf("%s%s/single/%s/", constant.UPLOAD_ROOT_DIR, upload.Model, upload.Single.VersionStr)
	dsds := fmt.Sprintf("%s%s/dsds/", constant.UPLOAD_ROOT_DIR, upload.Model)
	utils.TouchDir(single, 0)
	utils.TouchDir(dsds, 0)
	dsds_symlink := fmt.Sprintf("%s%s", dsds, upload.Dsds.VersionStr)
	err := os.Symlink(single, dsds_symlink)
	if err != nil {
		utils.Delete(single)
		utils.Delete(dsds_symlink)
		log.Log(tag, fmt.Sprintf("Failed to create symbal link: err = %s\n", err))
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Server internal error")
		return c.Redirect(Radio.PostUpload)
	}

	//parse information of each cp file
	file_list := make([]models.ImageFileComponent, len(utils.FILE_KEY_LIST))
	for index, key := range utils.FILE_KEY_LIST {
		log.Log(tag, fmt.Sprintf("check upload file : %s", key))
		fh_arr, ok := c.Params.Files[key]
		if !ok || len(fh_arr) <= 0 {
			utils.Delete(single)
			utils.Delete(dsds_symlink)
			log.Log(tag, fmt.Sprintf("loss upload file : %s", key))
			c.Validation.Keep()
			c.FlashParams()
			c.Flash.Error(fmt.Sprintf("loss upload file : %s", key))
			return c.Redirect(Radio.PostUpload)
		}
		fh := fh_arr[0]
		log.Log(tag, fmt.Sprintf("got upload file : %s", key))
		input, err := fh.Open()
		if err != nil {
			utils.Delete(single)
			utils.Delete(dsds_symlink)
			c.Validation.Keep()
			c.FlashParams()
			c.Flash.Error(fmt.Sprintf("Failed to open file %s :\n err = %s", fh.Filename, err))
			return c.Redirect(Radio.PostUpload)
		}
		defer input.Close()

		src := models.ImageFileComponent{}
		src.Name = fh.Filename
		src.Offset = utils.FILE_OFFSET_LIST[index]
		file_list[index] = src

		switch key {
		case "single_cp":
			upload.Single.Cp = src
		case "single_dsp":
			upload.Single.Dsp = src
		case "dsds_cp":
			upload.Dsds.Cp = src
		case "dsds_dsp":
			upload.Dsds.Dsp = src
		}

		err = utils.CopyFile(input, single+src.Name)
		if err != nil {
			utils.Delete(single)
			utils.Delete(dsds_symlink)
			c.Validation.Keep()
			c.FlashParams()
			c.Flash.Error("Server internal error")
			return c.Redirect(Radio.PostUpload)
		}
	}

	err = models.GenerateImageFile(single, file_list, constant.RADIO_IMAGE_SIZE, single+constant.RADIO_IMAGE_NAME)
	if err != nil {
		utils.Delete(single)
		utils.Delete(dsds_symlink)
		log.Log(tag, fmt.Sprintf("Failed to generate image file %s :\n err = %s", constant.UPLOAD_ROOT_DIR+constant.RADIO_IMAGE_NAME, err))
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Server internal error")
		return c.Redirect(Radio.PostUpload)
	}

	dal, err := models.NewDal(models.DRIVER, models.DNS)
	if err != nil {
		utils.Delete(single)
		utils.Delete(dsds_symlink)
		log.Log(tag, fmt.Sprintf("New dal error: %s\n", err))
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Server internal error")
		return c.Redirect(Radio.PostUpload)
	}
	defer dal.Close()
	err = upload.Save(dal)
	if err != nil {
		utils.Delete(single)
		utils.Delete(dsds_symlink)
		log.Log(tag, fmt.Sprintf("Failed to save upload info : %s,  err = %s\n", upload, err))
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Server internal error")
		return c.Redirect(Radio.PostUpload)
	}

	return c.Redirect(Radio.Index)
}
