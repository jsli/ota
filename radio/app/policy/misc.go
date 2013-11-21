package policy

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	cp_constant "github.com/jsli/cp_release/constant"
	"github.com/jsli/gtbox/file"
	"github.com/jsli/gtbox/pathutil"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"regexp"
	"strings"
	"time"
)

func RecordMd5(path, txt_path string) error {
	md5_str, err := file.Md5SumFile(path)
	if md5_str == "" || err != nil {
		return err
	} else {
		parent := pathutil.ParentPath(path)
		if parent == "" {
			return fmt.Errorf("record md5 failed: %s parent path is empty", path)
		}
		base_name := pathutil.BaseName(path)
		err := file.WriteString2File(fmt.Sprintf("%s %s", md5_str, base_name), txt_path)
		if err != nil {
			return err
		}
	}
	return nil
}

func Md5Dtim(dtim []byte) string {
	md5h := md5.New()
	md5h.Write(dtim)
	md5_str := hex.EncodeToString(md5h.Sum(nil))
	return md5_str
}

func FormatTime(time_unix int64) string {
	t := time.Unix(int64(time_unix), 0)
	return t.Format(ota_constant.TIME_FMT)
}

func ConvertModel(model string) string {
	model = strings.ToUpper(model)
	if strings.Contains(model, ota_constant.TAG_1088) ||
		strings.Contains(model, ota_constant.TAG_1T88) {
		model = ota_constant.MODEL_1088
	} else if strings.Contains(model, ota_constant.TAG_1920) {
		model = ota_constant.MODEL_1920
	}
	return model
}

func ConvertAndroidPlatform(platform string) string {
	if strings.HasPrefix(platform, ota_constant.TAG_JB_4_2) {
		platform = ota_constant.PLATFORM_JB_4_2
	} else if strings.HasPrefix(platform, ota_constant.TAG_JB_4_3) {
		platform = ota_constant.PLATFORM_JB_4_3
	}
	return platform
}

func FilterIp(ip string) string {
	return strings.Split(ip, ":")[0]
}

func GetFiltersFromFile(mode string, key string) []string {
	path := fmt.Sprintf("%s%s_%s", cp_constant.FILTER_ROOT, mode, key)
	content, err := file.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	content = strings.TrimSpace(content)
	return strings.Split(content, "\n")
}

func CheckImageByFilters(image string, filters []string) bool {
	for _, filter := range filters {
		if strings.Contains(image, filter) {
			return false
		}
	}
	return true
}

var VersionPattern = regexp.MustCompile(`\d+\.\d+\.\w{3}`)

func ReplaceVersionInPath(path string, version string) (string, error) {
	rep := fmt.Sprintf("${1}%s", version)
	r_path := VersionPattern.ReplaceAllString(path, rep)
	return r_path, nil
}
