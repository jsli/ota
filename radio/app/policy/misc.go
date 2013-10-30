package policy

import (
	"fmt"
	"github.com/jsli/gtbox/file"
	"github.com/jsli/gtbox/pathutil"
	ota_constant "github.com/jsli/ota/radio/app/constant"
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

func FormatTime(time_unix int64) string {
	t := time.Unix(int64(time_unix), 0)
	return t.Format(ota_constant.TIME_FMT)
}
