package policy

import (
	"fmt"
	"github.com/jsli/gtbox/file"
	"github.com/jsli/gtbox/pathutil"
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
