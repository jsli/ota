package constant

import ()

var (
	//ARBI|LTG|SINGLE|HLLTE/HLLTE_CP_2.29.000/Seagull/HL_LTG.bin
	//GRBI|LTG|SINGLE|HLLTE/HLLTE_CP_2.29.000/TTD_WK_NL_MSA_2.29.000/HL_DL_M09_Y0_AI_SKL_Flash.bin
	TestDataHLTD = [][]string{
		{"ARBI", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxx.bin"},
		{"GRBI", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxxFlash.bin"},
	}
	TestDataHLWB = [][]string{
		{"ARBI", "WG", "SINGLE", "HLWB/HLWB_CP_1.55.000/xxx.bin"},
		{"GRBI", "WG", "SINGLE", "HLWB/HLWB_CP_1.55.000/xxxFlash.bin"},
	}
	TestDataHLTD_DSDS = [][]string{
		{"ARBI", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxx.bin"},
		{"GRBI", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxxFlash.bin"},
		{"ARB2", "TG", "DSDS", "HLTD_DSDS/HLTD_DSDS_CP_3.28.000/xxx.bin"},
		{"GRB2", "TG", "DSDS", "HLTD_DSDS/HLTD_DSDS_CP_3.28.000/xxxFlash.bin"},
	}
	TestDataHLWB_DSDS = [][]string{
		{"ARBI", "WG", "SINGLE", "HLWB/HLWB_CP_1.55.000/xxx.bin"},
		{"GRBI", "WG", "SINGLE", "HLWB/HLWB_CP_1.55.000/xxxFlash.bin"},
		{"ARB2", "WG", "DSDS", "HLWB_DSDS/HLWB_CP_2.58.917/xxx.bin"},
		{"GRB2", "WG", "DSDS", "HLWB_DSDS/HLWB_CP_2.58.917/xxxFlash.bin"},
	}
	TestDataHLTDR = [][]string{
		{"ARBI", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxx.bin"},
		{"GRBI", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxxFlash.bin"},
		{"RFIC", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxxrfic.bin"},
	}
	TestDataLTER = [][]string{
		{"ARBI", "LWG", "SINGLE", "LWG/HL_CP_2.30.000/HL_CP/Seagull/HL_LWG_DKB.bin"},
		{"GRBI", "LWG", "SINGLE", "LWG/HL_CP_2.30.000/HL_MSA_2.30.000/HL_LWG_M09_B0_SKL_Flash.bin"},
		{"RFIC", "LWG", "SINGLE", "LWG/HL_CP_2.30.000/RFIC/1920_FF/Skylark_LWG.bin"},
		{"ARB2", "LTG", "SINGLE", "LTG/HL_CP_3.30.000/HL_CP/Seagull/HL_LTG_DL_DKB.bin"},
		{"GRB2", "LTG", "SINGLE", "LTG/HL_CP_3.30.000/HL_MSA_3.30.000/HL_DL_M09_Y0_AI_SKL_Flash.bin"},
		{"RFI2", "LTG", "SINGLE", "LTG/HL_CP_3.30.000/RFIC/1920_FF/Skylark_LTG.bin"},
	}
)
