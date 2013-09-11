package models

import (
	"fmt"
	"github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/utils"
	"github.com/jsli/ota/radio/app/log"
	"github.com/robfig/revel"
	"io"
	"os"
)

type ImageFileComponent struct {
	Name   string
	Offset int64
}

func (comp ImageFileComponent) String() string {
	return fmt.Sprintf("ImageFileComponent(Name=%s, Offset=%d)",
		comp.Name, comp.Offset)
}

type CpFile struct {
	Type       int `1:single, 2:dsds`
	VersionInt int64
	VersionStr string
	Cp         ImageFileComponent
	Dsp        ImageFileComponent
}

func (cp CpFile) String() string {
	return fmt.Sprintf("CpFile(Type=%d, v_int=%d, v_str=%s, cp=%s, dsp=%s)",
		cp.Type, cp.VersionInt, cp.VersionStr, cp.Cp, cp.Dsp)
}

type RadioImageFile struct {
	Model  string
	Single CpFile
	Dsds   CpFile
}

func (f RadioImageFile) String() string {
	return fmt.Sprintf("RadioImageFile(model=%s, Single=%s, Dsds=%s)",
		f.Model, f.Single, f.Dsds)
}

func (f *RadioImageFile) Validate(v *revel.Validation) {
	v.Check(f.Model,
		revel.Required{},
		LegalModelValidator{},
	).Message("Illegal model")

	v.Check(f.Single.VersionStr,
		revel.Required{},
		revel.Match{constant.CP_VERSION_REX},
	).Message("Illegal single version")
	
	v.Check(f.Dsds.VersionStr,
		revel.Required{},
		revel.Match{constant.CP_VERSION_REX},
	).Message("Illegal dsds version")
}

type LegalModelValidator struct {
}

func (legal LegalModelValidator) IsSatisfied(obj interface{}) bool {
	return utils.IsAvailableModel(obj.(string))
}

func (legal LegalModelValidator) DefaultMessage() string {
	return "Illegal model"
}

func (f *RadioImageFile) Save(dal *Dal) error {
	tag := "RadioImageFile.Save"
	stmt, err := dal.Link.Prepare("INSERT radio_image SET model=?,single_version_str=?,single_version=?,dsds_version_str=?,dsds_version=?")
	if err != nil {
		log.Log(tag, fmt.Sprintf("Prepare error : %s", err))
		return err
	}
	_, err = stmt.Exec(f.Model, f.Single.VersionStr, f.Single.VersionInt, f.Dsds.VersionStr, f.Dsds.VersionInt)
	if err != nil {
		log.Log(tag, fmt.Sprintf("Exec error : %s", err))
		return err
	}

	//	id, err := res.LastInsertId()
	return nil
}

func GenerateImageFile(comp_file_root string, comp_list []ImageFileComponent, img_size int64, dest string) error {
	tag := "GenerateImageFile"
	img_file, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, constant.IMAGE_FILE_ACCESS)
	if err != nil {
		log.Log(tag, fmt.Sprintf("Failed to open file %s :\n err = %s", dest, err))
		return err
	}
	defer img_file.Close()

	buffer := make([]byte, 4096)

	for _, src := range comp_list {
		f, err := os.OpenFile(comp_file_root+src.Name, os.O_RDONLY, 0)
		if err != nil {
			log.Log(tag, fmt.Sprintf("Failed to open file %s :\n err = %s", src.Name, err))
			return err
		}
		defer f.Close() // dup close opt, in case abnormally return in the for-cycle below

		img_file.Seek(src.Offset, 0)
		for {
			n, err := f.Read(buffer)
			if n > 0 && err == nil {
				n, err = img_file.Write(buffer[:n])
				if err != nil {
					log.Log(tag, fmt.Sprintf("Failed to write file %s :\n err = %s", dest, err))
					return err
				}
				continue
			} else if err == io.EOF {
				break
			} else if err != nil {
				log.Log(tag, fmt.Sprintf("Failed to read file %s :\n err = %s", src.Name, err))
				return err
			}
		}
		f.Close()
	}
	img_file.Truncate(img_size)
	log.Log(tag, fmt.Sprintf("Successfully to generate image file %s \n", dest))
	return nil
}
