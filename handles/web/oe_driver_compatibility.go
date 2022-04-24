package web

import (
	"cve-sa-backend/dao"
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
)

type DriverHandle struct {
}

func (d *DriverHandle) FindAllDriverCompatibility(req cveSa.OeCompSearchRequest) (*cveSa.ResultData, error) {
	driverCompatibility, total, err := dao.DefaultCompatibilityDriver.FindAllDriver(req)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}

	return d.returnDriverCompatibility(driverCompatibility, total), nil
}

func (d *DriverHandle) returnDriverCompatibility(datas []models.OeCompatibilityDriver, total int64) *cveSa.ResultData {
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

func (d *DriverHandle) GetOsList(lang string) (data []string, err error) {
	data, err = dao.DefaultCompatibilityDriver.GetOsList(lang)
	return
}

func (d *DriverHandle) GetArchitectureList(lang string) (data []string, err error) {
	data, err = dao.DefaultCompatibilityDriver.GetArchitectureList(lang)
	return
}
