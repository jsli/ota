package log

import (
	"github.com/jsli/ota/radio/app/constant"
	"os"
	"fmt"
	"time"
)

func Log(tag, msg string) {
	//fmt: <time> -- <tag> : <message>
	log := fmt.Sprintf("<%s> -- <%s> : %s\n", time.Now(), tag, msg)
	writeString2File(log, constant.LOG_ROOG_DIR+constant.UPLOAD_LOG_NAME)
	fmt.Printf(log)
}

func writeString2File(content, path string) error {
	destFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, constant.DEFAULT_FILE_ACCESS)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
