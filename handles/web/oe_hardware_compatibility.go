package web

import (
	"cve-sa-backend/dao"
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
)

type HardwareHandle struct {
}

func (h *HardwareHandle) FindAllHardwareCompatibility(req cveSa.OeCompSearchRequest) (*cveSa.ResultData, error) {
	hardware, total, err := dao.DefaultCompatibilityHardware.FindAllHardware(req)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}

	return h.returnHardware(hardware, total), nil
}

func (h *HardwareHandle) returnHardware(datas []*models.OeCompatibilityHardware, total int64) *cveSa.ResultData {
	var hardwareCompatibility = make([]cveSa.HardwareCompatibility, 0, len(datas))
	for _, v := range datas {
		hardwareCompatibility = append(hardwareCompatibility, cveSa.HardwareCompatibility{
			OeCompatibilityHardware: *v,
			BoardCards:              make([]models.OeCompatibilityHardwareAdapter, 0)})
	}
	return &cveSa.ResultData{SecurityNoticeList: make([]cveSa.SecurityNoticeData, 0), CveDatabaseList: make([]cveSa.DatabaseData, 0),
		ApplicationCompList: make([]models.ROeCompatibilityApplication, 0), HardwareCompList: hardwareCompatibility,
		DriverCompList: make([]models.OeCompatibilityDriver, 0), Total: int(total),
	}
}

func (h *HardwareHandle) GetHardwareCompatibilityById(id int64) (*cveSa.HardwareCompatibility, error) {
	data, err := dao.DefaultCompatibilityHardware.GetOneHardware(&models.OeCompatibilityHardware{Id: id})
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	return &cveSa.HardwareCompatibility{OeCompatibilityHardware: *data, BoardCards: make([]models.OeCompatibilityHardwareAdapter, 0)}, nil
}

func (h *HardwareHandle) GetCpuList(lang string) (datas []string, err error) {
	datas, err = dao.DefaultCompatibilityHardware.GetCpuList(lang)
	return
}

func (h *HardwareHandle) ByhardwareId(hardwareId int64) (datas []models.OeCompatibilityHardwareAdapter, err error) {
	datas, err = dao.DefaultCompatibilityHardwareAdapter.ByhardwareId(hardwareId)
	return
}

func (h *HardwareHandle) GetOsForHardware(lang string) (datas []string, err error) {
	datas, err = dao.DefaultCompatibilityHardware.GetOsForHardware(lang)
	return
}

func (h *HardwareHandle) GetArchitectureListForHardware(lang string) (datas []string, err error) {
	datas, err = dao.DefaultCompatibilityHardware.GetArchitectureListForHardware(lang)
	return
}
