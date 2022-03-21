package main

import (
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"cve-sa-backend/utils"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/parsexml"
)

func Test_main(t *testing.T) {
	cveFileName := "2022/cvrf-openEuler-SA-2022-1495.xml"
	fileByte, err := utils.GetCvrfFile(cveFileName)
	if err != nil {
		fmt.Println(err)
	}

	Element := utils.FixedCveXml{}
	err = xml.Unmarshal(fileByte, &Element)
	if err != nil {
		fmt.Println(err)
	}

	var list []cveSa.DatabaseData
	updateTime := time.Now()

	for _, v := range Element.Vulnerability {

		cve, err := parsexml.GetCVEDatabase("", "", v, updateTime)
		if err != nil {
			fmt.Println(err)
		}
		list = append(list, cve)
	}
}
