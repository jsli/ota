package models

import (
	"fmt"
)

type DeviceInfo struct {
	Model    string //refer to constant MODEL_LIST
	Platform string
	MacAddr  string
}

func (di DeviceInfo) String() string {
	return fmt.Sprintf("DeviceInfo( Model=%s, MacAddr=%s, Platform=%s)", di.Model, di.MacAddr, di.Platform)
}
