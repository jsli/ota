package models

import (
	"fmt"
)

type DeviceInfo struct {
	Model   string //refer to constant MODEL_LIST
	MacAddr string
}

func (di DeviceInfo) String() string {
	return fmt.Sprintf("DeviceInfo( Model=%s, MacAddr=%s)", di.Model, di.MacAddr)
}
