package models

import (
	"fmt"
)

type UpdateRequest struct {
	Device DeviceInfo
	Cps    []CpRequest
}

func (ur UpdateRequest) String() string {
	return fmt.Sprintf("UpdateRequest( Device=%s, Cps=%s)", ur.Device, ur.Cps)
}

type CpRequest struct {
	Mode    string
	Version string
	Images  map[string]string
}

func (cr CpRequest) String() string {
	return fmt.Sprintf("CpRequest( Mode=%s, Version=%s, Images=%s)", cr.Mode, cr.Version, cr.Images)
}
