package models

import ()

/*
 * for list-style result
 */
type ModeNode struct {
	Mode    string   `json:"mode_name"`
	CpArray []CpNode `json:"cp_array"`
}

type CpNode struct {
	VersionNo  string      `json:"version_no"`
	ImageArray []ImageNode `json:"image_array"`
}

type ImageNode struct {
	ImageName string   `json:"image_name"`
	Images    []string `json:"images"`
}

/*
 * for current cp information
 */
type CpAndImages struct {
	Version string `json:"version"`

	/*image-id as key*/
	Images map[string]string `json:"images"`
}

/* mode as key */
type CurrentCps map[string]CpAndImages
