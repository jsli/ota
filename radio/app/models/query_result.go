package models

import (
	ota_constant "github.com/jsli/ota/radio/app/constant"
)

/*******************************
HLTD: {
	2.59.000: {
		ARBEL: [1]
			0:  "HL/HLTD/HLTD_CP_2.59.000/Seagull/HL_TD_CP.bin"
		MSA: [0]
			0:  "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
			1:  "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}
*******************************/
type ImageList []string

/* image's id as key, like ARBI, GRBI... */
type ImagesList map[string]ImageList

/* version as key */
type AvailableCpComponent map[string]ImagesList

/* mode as key */
type AvailableCps map[string]AvailableCpComponent

/*******************************
images: {
	ARBEL: "HL/HLTD/HLTD_CP_2.51.000/Seagull/HL_TD_CP.bin"
	MSA: "HL/HLTD/HLTD_CP_2.51.000/HLTD_MSA_2.51.000/A0/HL_TD_M08_AI_A0_Flash.bin"
}
 *******************************/
/* image's id as key, like ARBI, GRBI... */
type Images map[string]string

/*******************************
HLTD: {
	images: {...}-
	version: "2.51.000"
}
 *******************************/
type CurrentCpComponent struct {
	Version string
	Images  Images
}

/* mode as key */
type CurrentCps map[string]CurrentCpComponent

type QueryData struct {
	Available AvailableCps `json:"available"`
	Current   CurrentCps   `json:"current"`
}

type QueryResult struct {
	Data QueryData `json:"data"`
	Result
}

func NewQueryResult() *QueryResult {
	result := QueryResult{}
	result.Data = QueryData{}
	result.Extra = ExtraInfo{ApiVersion: API_VERSION, ErrorCode: ota_constant.ERROR_CODE_NOERR, ErrorMessage: nil}
	return &result
}
