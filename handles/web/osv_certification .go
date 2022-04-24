package web

import (
	"encoding/json"

	"cve-sa-backend/dao"
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	_const "cve-sa-backend/utils/const"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
)

type OsvHandle struct {
}

func (o *OsvHandle) FindAllOsv(req cveSa.RequestOsv) (*cveSa.ResultOSVData, error) {
	datas, total, err := dao.DefaultCompatibilityOsv.OSVFindAll(req)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}

	OSVList := make([]models.ROeCompatibilityOsv, 0, len(datas))
	for _, v := range datas {
		var t []models.Record
		_ = json.Unmarshal([]byte(v.ToolsResult), &t)
		var p []models.Record
		_ = json.Unmarshal([]byte(v.PlatformResult), &p)

		OSVList = append(OSVList, models.ROeCompatibilityOsv{
			OeCompatibilityOsv: v,
			ToolsResult:        t,
			PlatformResult:     p,
			Updateime:          v.Updateime.Format(_const.Format),
		})
	}

	result := &cveSa.ResultOSVData{
		Total:   int(total),
		OsvList: OSVList,
	}

	return result, nil
}

func (o *OsvHandle) GetOsvName() (data []string, err error) {
	return dao.DefaultCompatibilityOsv.GetOsvName()
}

func (o *OsvHandle) GetType() (data []string, err error) {
	return dao.DefaultCompatibilityOsv.GetType()
}

func (o *OsvHandle) GetOne(id int64) (*models.ROeCompatibilityOsv, error) {

	osv, err := dao.DefaultCompatibilityOsv.GetOneOSV(&models.OeCompatibilityOsv{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	var t []models.Record
	_ = json.Unmarshal([]byte(osv.ToolsResult), &t)
	var p []models.Record
	_ = json.Unmarshal([]byte(osv.PlatformResult), &p)

	return &models.ROeCompatibilityOsv{
		OeCompatibilityOsv: *osv,
		ToolsResult:        t,
		PlatformResult:     p,
		Updateime:          osv.Updateime.Format(_const.Format),
	}, nil
}
