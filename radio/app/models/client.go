package models

import (
	"fmt"
)

type DeviceInfo struct {
	Model    string //refer to constant MODEL_LIST
	Platform string
}

func (di DeviceInfo) String() string {
	return fmt.Sprintf("DeviceInfo( Model=%s, Platform=%s)", di.Model, di.Platform)
}
