package web

import (
	"cve-sa-backend/dao"
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
)

func FindAllHardwareCompatibility(req cveSa.OeCompSearchRequest) (*cveSa.ResultData, error) {
	hardware, total, err := dao.FindAllHardware(req)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}

	return returnHardware(hardware, total), nil
}

func returnHardware(datas []*models.OeCompatibilityHardware, total int64) *cveSa.ResultData {
	var hardwareCompatibility = make([]cveSa.HardwareCompatibility, 0, len(datas))
	for _, v := range datas {
		hardwareCompatibility = append(hardwareCompatibility, cveSa.HardwareCompatibility{
			OeCompatibilityHardware: v,
			BoardCards:              make([]models.OeCompatibilityHardwareAdapter, 0)})
	}
	return &cveSa.ResultData{SecurityNoticeList: make([]cveSa.SecurityNoticeData, 0), CveDatabaseList: make([]cveSa.DatabaseData, 0),
		ApplicationCompList: make([]models.ROeCompatibilityApplication, 0), HardwareCompList: hardwareCompatibility,
		DriverCompList: make([]models.OeCompatibilityDriver, 0), Total: int(total),
	}
}

func GetHardwareCompatibilityById(id int64) (*cveSa.HardwareCompatibility, error) {
	data, err := dao.GetOneHardware(&models.OeCompatibilityHardware{Id: id})
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	return &cveSa.HardwareCompatibility{OeCompatibilityHardware: data, BoardCards: make([]models.OeCompatibilityHardwareAdapter, 0)}, nil
}

func GetCpuList(lang string) (datas []string, err error) {
	datas, err = dao.GetCpuList(lang)
	return
}

func ByhardwareId(hardwareId int64) (datas []models.OeCompatibilityHardwareAdapter, err error) {
	datas, err = dao.ByhardwareId(hardwareId)
	return
}

func GetOsForHardware(lang string) (datas []string, err error) {
	datas, err = dao.GetOsForHardware(lang)
	return
}

func GetArchitectureListForHardware(lang string) (datas []string, err error) {
	datas, err = dao.GetArchitectureListForHardware(lang)
	return
}
