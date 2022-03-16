package web

import (
	"cve-sa-backend/dao"
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
)

func FindAllDriverCompatibility(req cveSa.OeCompSearchRequest) (*cveSa.ResultData, error) {
	driverCompatibility, total, err := dao.FindAllDriver(req)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}

	return returnDriverCompatibility(driverCompatibility, total), nil
}

func returnDriverCompatibility(datas []models.OeCompatibilityDriver, total int64) *cveSa.ResultData {
	if len(datas) == 0 {
		datas = make([]models.OeCompatibilityDriver, 0)
	}
	return &cveSa.ResultData{
		SecurityNoticeList:  make([]cveSa.SecurityNoticeData, 0),
		CveDatabaseList:     make([]cveSa.DatabaseData, 0),
		ApplicationCompList: make([]models.ROeCompatibilityApplication, 0),
		HardwareCompList:    make([]cveSa.HardwareCompatibility, 0),
		DriverCompList:      datas,
		Total:               int(total),
	}
}

func GetOsList(lang string) (data []string, err error) {
	data, err = dao.GetOsList(lang)
	return
}

func GetArchitectureList(lang string) (data []string, err error) {
	data, err = dao.GetArchitectureList(lang)
	return
}
