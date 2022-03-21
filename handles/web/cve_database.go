package web

import (
	"cve-sa-backend/dao"
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	_const "cve-sa-backend/utils/const"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
)

func FindAllCVEDatabase(req cveSa.RequestData) (*cveSa.ResultData, error) {
	datas, total, err := dao.DatabaseFindAll(req)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}
	return ReturnCVEDatabase(datas, total), nil
}

func ReturnCVEDatabase(data []models.CveDatabase, total int64) *cveSa.ResultData {
	cveDatabaseList := make([]cveSa.DatabaseData, 0, len(data))
	for _, v := range data {
		rc := models.RCveDatabase{
			CveDatabase: v,
			Updateime:   v.Updateime.Format(_const.Format),
		}
		cveDatabaseList = append(cveDatabaseList, cveSa.DatabaseData{
			RCveDatabase: rc,
		})
	}

	return &cveSa.ResultData{
		SecurityNoticeList:  make([]cveSa.SecurityNoticeData, 0),
		CveDatabaseList:     cveDatabaseList,
		ApplicationCompList: make([]models.ROeCompatibilityApplication, 0),
		HardwareCompList:    make([]cveSa.HardwareCompatibility, 0),
		DriverCompList:      make([]models.OeCompatibilityDriver, 0),
		Total:               int(total),
	}
}
