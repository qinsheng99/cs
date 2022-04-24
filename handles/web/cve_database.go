package web

import (
	"cve-sa-backend/dao"
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	_const "cve-sa-backend/utils/const"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
)

type CveDatabaseHandle struct {
}

func (cv *CveDatabaseHandle) FindAllCVEDatabase(req cveSa.RequestData) (*cveSa.ResultData, error) {
	datas, total, err := dao.DefaultCveDatabase.DatabaseFindAll(req)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}
	return cv.ReturnCVEDatabase(datas, total), nil
}

func (cv *CveDatabaseHandle) ReturnCVEDatabase(data []models.CveDatabase, total int64) *cveSa.ResultData {
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

func (cv *CveDatabaseHandle) GetByCveIdAndPackageName(cveId, packageName string) (*cveSa.DatabaseData, error) {
	result := &cveSa.DatabaseData{}
	cve, err := dao.DefaultCveDatabase.GetOneDatabaseTypeTwo(&models.CveDatabase{
		CveId:       cveId,
		PackageName: packageName,
	})
	if utils.ErrorNotFound(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	result.RCveDatabase = models.RCveDatabase{
		CveDatabase: *cve,
		Updateime:   cve.Updateime.Format(_const.Format),
	}

	list, err := dao.DefaultSecurityNotice.NoticeByCveIdAndAffectedComponent(cve.CveId, cve.PackageName)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		result.SecurityNoticeNo = list[0].SecurityNoticeNo
	}

	return result, nil
}

func (cv *CveDatabaseHandle) GetCVEProductPackageListByCveId(cveId string) ([]models.RCveProductPackage, error) {
	var result []models.RCveProductPackage
	list, _, err := dao.DefaultCveProductPackage.GetProductPackageListTypeTwo(&models.CveProductPackage{
		CveId: cveId,
	})

	if err != nil {
		return result, err
	}

	for _, v := range list {
		result = append(result, models.RCveProductPackage{
			CveProductPackage: v,
			Updateime:         v.Updateime.Format(_const.Format),
		})
	}

	return result, nil

}

func (cv *CveDatabaseHandle) GetCVEProductPackageList(cveId, packageName string) ([]models.RCveProductPackage, error) {
	var result []models.RCveProductPackage
	list, _, err := dao.DefaultCveProductPackage.GetProductPackageListTypeTwo(&models.CveProductPackage{
		CveId:       cveId,
		PackageName: packageName,
	})

	if err != nil {
		return result, err
	}

	for _, v := range list {
		result = append(result, models.RCveProductPackage{
			CveProductPackage: v,
			Updateime:         v.Updateime.Format(_const.Format),
		})
	}

	return result, nil
}
