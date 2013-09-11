package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/log"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	SLASH = string(os.PathSeparator)
)

/*
 *path must end with "/"
 */
func IsDirExist(path string) (bool, error) {
	tag := "IsDirExist"
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		log.Log(tag, fmt.Sprintf("IsDirExist info: %s stat error", path))
		return false, err
	}
	return true, nil
}

func TouchDir(path string, mode os.FileMode) error {
	exist, err := IsDirExist(path)
	if mode == 0 {
		mode = constant.DEFAULT_DIR_ACCESS
	}
	if !exist && err == nil {
		return os.MkdirAll(path, mode)
	} else if exist && err == nil {
		return nil
	} else {
		return err
	}
}

func ExtractZipFile(name string, dest string) error {
	tag := "ExtractZipFile"
	rc, err := zip.OpenReader(name)
	if err != nil {
		log.Log(tag, fmt.Sprintf("%s OpenReader failed : %s\n", name, err))
		return err
	}
	defer rc.Close()

	for _, zf := range rc.File {
		if strings.HasSuffix(zf.Name, SLASH) {
			os.MkdirAll(dest+SLASH+zf.Name, zf.Mode())
			continue
		}

		reader, err := zf.Open()
		if err != nil {
			log.Log(tag, fmt.Sprintf("%s Open failed : %s\n", zf.Name, err))
			return err
		}

		//There is a BUG!!!S
		//can't save ModeSymlink to it !!!!!!
		fw, err := os.OpenFile(dest+SLASH+zf.Name, os.O_CREATE|os.O_WRONLY, zf.Mode())
		if err != nil {
			log.Log(tag, fmt.Sprintf("%s Open failed : %s\n", dest+SLASH+zf.Name, err))
			return err
		}

		if _, err := io.Copy(fw, reader); err != nil {
			log.Log(tag, fmt.Sprintf("%s Copy failed : %s\n", dest+SLASH+zf.Name, err))
			return err
		}
		fw.Close()
		reader.Close()
	}
	return nil
}

func MakeZipFile(src, dest string) error {
	tag := "MakeZipFile"
	list := ListFilesRecursive("", src, false)
	zw, err := os.Create(dest)
	if err != nil {
		log.Log(tag, fmt.Sprintf("%s create zip dest failed : %s\n", dest, err))
		return err
	}
	defer zw.Close()

	w := zip.NewWriter(zw)
	if w == nil {
		log.Log(tag, fmt.Sprintf("%s create zip dest writer failed\n", dest))
		return err
	}
	defer w.Close()

	for _, str := range list {
		f, err := w.Create(str)
		if err != nil {
			log.Log(tag, fmt.Sprintf("create %s failed\n", str))
			return err
		}
		reader, err := os.OpenFile(src+str, os.O_RDONLY, 0644)
		if err != nil {
			log.Log(tag, fmt.Sprintf("%s open file failed : %s\n", src+str, err))
			return err
		}

		if _, err := io.Copy(f, reader); err != nil {
			log.Log(tag, fmt.Sprintf("%s copy file failed : %s\n", src+str, err))
			return err
		}
		reader.Close()
	}
	return nil
}

func ListFilesRecursive(prefix, path string, b bool) []string {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}

	list := make([]string, 0, 10)
	var dir_name string
	if !b {
		dir_name = ""
	} else {
		dir_name = BaseName(path) + SLASH
	}
	for _, info := range fileInfos {
		if info.IsDir() {
			tmp_list := ListFilesRecursive(prefix+dir_name, path+info.Name()+SLASH, true)
			list = append(list, tmp_list...)
		} else if info.Mode().IsRegular() {
			list = append(list, prefix+dir_name+info.Name())
		}
	}
	return list
}

func BaseName(path string) string {
	list := SplitPath(path)
	if strings.HasSuffix(path, SLASH) {
		list = list[:len(list)-1]
	}
	if list != nil && len(list) > 0 {
		return list[len(list)-1]
	}
	return ""
}

func SplitPath(path string) []string {
	return strings.Split(path, SLASH)
}

func ParentPath(path string) string {
	list := SplitPath(path)
	var isAbs bool = false
	//	var isSlashEnd bool = false
	if strings.HasPrefix(path, SLASH) {
		list = list[1:]
		isAbs = true
	}

	if strings.HasSuffix(path, SLASH) {
		list = list[:len(list)-1]
		//		isSlashEnd = true
	}

	if len(list) <= 0 {
		return SLASH
	} else {
		list = list[:len(list)-1]
		if len(list) <= 0 {
			if isAbs {
				return SLASH
			}
			return ""
		} else {
			parent := strings.Join(list, SLASH)
			if isAbs {
				parent = SLASH + parent
			}
			if parent != "" {
				parent += SLASH
			}
			return parent
		}
	}
}

func CopyFileWithPath(src, dest string) (int64, error) {
	tag := "CopyFileWithPath"
	fi, fi_err := os.Stat(src)
	if fi_err != nil {
		log.Log(tag, fmt.Sprintf("%s stat error : %s", src, fi_err))
		return 0, fi_err
	}
	srcFile, err := os.OpenFile(src, os.O_RDONLY, fi.Mode())
	if err != nil {
		log.Log(tag, fmt.Sprintf("%s open file error : %s", src, err))
		return 0, err
	}
	defer srcFile.Close()

	parent := ParentPath(dest)
	TouchDir(parent, 0755)
	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, fi.Mode())
	if err != nil {
		log.Log(tag, fmt.Sprintf("%s open file error : %s", dest, err))
		return 0, err
	}
	defer destFile.Close()

	n, err := io.Copy(destFile, srcFile)
	return n, err
}

/*
 *src and dest must be end with "/"
 */
func CopyDirWithPath(src, dest string) (bool, error) {
	tag := "CopyDirWithPath"
	mode, err := GetFileMode(src)
	if err != nil {
		log.Log(tag, fmt.Sprintf("%s get file mode error : %s", src, err))
		return false, err
	}
	dir_name := BaseName(src)
	err = TouchDir(dest+dir_name+SLASH, mode)
	if err != nil {
		log.Log(tag, fmt.Sprintf("CopyDir touch dir error : %s", err))
		return false, err
	}
	fileInfos, err := ioutil.ReadDir(src)
	if err != nil {
		log.Log(tag, fmt.Sprintf("CopyDir read dir error : %s", err))
		return false, err
	}
	for _, info := range fileInfos {
		if info.IsDir() {
			res, err := CopyDirWithPath(src+info.Name()+SLASH, dest+dir_name+SLASH)
			if !res || err != nil {
				log.Log(tag, fmt.Sprintf("copy dir error : %s", err))
				return false, errors.New(fmt.Sprintf("copy dir error : %s", err))
			}
		} else if info.Mode().IsRegular() {
			_, err := CopyFileWithPath(src+info.Name(), dest+dir_name+SLASH+info.Name())
			if err != nil {
				log.Log(tag, fmt.Sprintf("copy file error : %s", err))
				return false, errors.New(fmt.Sprintf("copy file error : %s", err))
			}
		}
	}
	return true, nil
}

func GetFileMode(path string) (os.FileMode, error) {
	tag := "GetFileMode"
	fi, err := os.Stat(path)
	if err != nil {
		log.Log(tag, fmt.Sprintf("%s stat error : %s", path, err))
		return 0, err
	}
	return fi.Mode(), nil
}

func CopyFile(src io.Reader, dest string) error {
	tag := "CopyFile"
	content, err := ioutil.ReadAll(src)
	if err != nil || len(content) == 0 {
		log.Log(tag, fmt.Sprintf("Failed to read file:\n err = %s", err))
		return err
	}

	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, constant.DEFAULT_FILE_ACCESS)
	if err != nil {
		log.Log(tag, fmt.Sprintf("Failed to open file %s :\n err = %s", dest, err))
		return err
	}
	defer destFile.Close()

	_, err = destFile.Write(content)
	if err != nil {
		log.Log(tag, fmt.Sprintf("Failed to write file %s :\n err = %s", dest, err))
		return err
	}
	log.Log(tag, fmt.Sprintf("Copy to %s success!", dest))
	return nil
}

func GenerateOtaPackage(cmd string, params []string) bool {
	tag := "GenerateOtaPackage"
	res := ExecCmd(cmd, params)
	if !res {
		log.Log(tag, fmt.Sprintf("%s exec failed: %s", res, params))
		return false
	}
	return true
}

func ExecCmd(cmd_str string, params []string) bool {
	tag := "ExecCmd"
	cmd := exec.Command(cmd_str, params...)
	buf, err := cmd.Output()
	if err != nil {
		log.Log(tag, fmt.Sprintf("The command failed to perform: %s", err))
		return false
	}
	log.Log(tag, fmt.Sprintf("Result: %s", buf))
	return true
}

func Delete(path string) error {
	tag := "Delete"
	err := os.RemoveAll(path)
	log.Log(tag, fmt.Sprintf("delete %s, result : %s", path, err))
	return err
}
