package web

import (
	"strings"

	"cve-sa-backend/dao"
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	_const "cve-sa-backend/utils/const"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
)

func FindAllSecurity(req cveSa.RequestData) (*cveSa.ResultData, error) {
	datas, total, err := dao.SecurityFindAll(req)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}
	return returnNoticeData(datas, total), nil
}

func returnNoticeData(datas []models.CveSecurityNotice, total int64) *cveSa.ResultData {
	securityNoticeData := SecurityNoticeData(datas)
	return &cveSa.ResultData{
		SecurityNoticeList:  securityNoticeData,
		CveDatabaseList:     make([]cveSa.DatabaseData, 0),
		ApplicationCompList: make([]models.ROeCompatibilityApplication, 0),
		HardwareCompList:    make([]cveSa.HardwareCompatibility, 0),
		DriverCompList:      make([]models.OeCompatibilityDriver, 0),
		Total:               int(total),
	}
}

func GetSecurityNoticePackageByPackageName(pname string) ([]models.RCveSecurityNoticePackage, error) {
	pname = utils.TrimString(pname)
	pnames := strings.Split(pname, ",")

	datas, err := dao.ByPackageName(pnames)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}
	return packageData(datas), nil
}

func packageData(datas []models.CveSecurityNoticePackage) []models.RCveSecurityNoticePackage {
	var list = make([]models.RCveSecurityNoticePackage, 0, len(datas))
	for _, v := range datas {
		list = append(list, models.RCveSecurityNoticePackage{CveSecurityNoticePackage: v, Updateime: v.Updateime.Format(_const.Format)})
	}
	return list
}

func NoticeByCVEID(cveId string) ([]cveSa.SecurityNoticeData, error) {
	datas, err := dao.NoticeByCveId(cveId)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}
	return SecurityNoticeData(datas), nil
}

func SecurityNoticeData(datas []models.CveSecurityNotice) []cveSa.SecurityNoticeData {
	var securityNoticeData = make([]cveSa.SecurityNoticeData, 0, len(datas))

	for _, v := range datas {
		c := cveSa.SecurityNoticeData{
			RCveSecurityNotice: &models.RCveSecurityNotice{
				CveSecurityNotice: v,
				Updateime:         v.Updateime.Format(_const.Format),
			},
			PackageHelperList: make([]cveSa.SAPackageHelper, 0),
			PackageList:       make([]models.RCveSecurityNoticePackage, 0),
			ReferenceList:     make([]models.RCveSecurityNoticeReference, 0),
			CveList:           make([]cveSa.DatabaseData, 0),
		}
		securityNoticeData = append(securityNoticeData, c)
	}
	return securityNoticeData
}

func SecurityNoticeDataOne(datas models.CveSecurityNotice) cveSa.SecurityNoticeData {
	return cveSa.SecurityNoticeData{
		RCveSecurityNotice: &models.RCveSecurityNotice{
			CveSecurityNotice: datas,
			Updateime:         datas.Updateime.Format(_const.Format),
		},
		PackageHelperList: make([]cveSa.SAPackageHelper, 0),
		PackageList:       make([]models.RCveSecurityNoticePackage, 0),
		ReferenceList:     make([]models.RCveSecurityNoticeReference, 0),
		CveList:           make([]cveSa.DatabaseData, 0),
	}
}

func ByCveIdAndAffectedComponent(cveId, affectedComponent string) ([]cveSa.SecurityNoticeData, error) {
	if cveId == "" {
		return nil, nil
	}

	datas, err := dao.NoticeByCveIdComponent(cveId, affectedComponent)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}
	return SecurityNoticeData(datas), nil
}

func NoticeBySecurityNoticeNo(s string) (*cveSa.SecurityNoticeData, error) {
	var SAPackages = make([]cveSa.SAPackageHelper, 0)
	securityNotice, err := dao.NoticeByNo(s)
	if err != nil {
		return nil, err
	}
	if securityNotice == nil {
		return nil, nil
	}
	snData := SecurityNoticeDataOne(*securityNotice)
	if securityNotice.AffectedProduct != "" {
		products := strings.Split(securityNotice.AffectedProduct, ";")
		for _, v := range products {
			packages, err := dao.NoticePackageByNoProduct(s, v)
			if err != nil {
				iniconf.SLog.Error(err)
				return &snData, nil
			}
			var SAPackage cveSa.SAPackageHelper
			SAPackage.ProductName = v
			SAPackage.Child = getSAPackageHelper(packages)
			SAPackages = append(SAPackages, SAPackage)
		}
	}
	snData.PackageHelperList = SAPackages
	references, err := dao.GetReferenceByNo(s)
	if err != nil {
		iniconf.SLog.Error(err)
		return &snData, nil
	}
	snData.ReferenceList = reReference(references)
	return &snData, nil
}

func getSAPackageHelper(datas []models.CveSecurityNoticePackage) []cveSa.SAPackageHelper {
	var SAPackageMap = make(map[string][]cveSa.SAPackageHelper)
	var SAPackages = make([]cveSa.SAPackageHelper, 0)
	if len(datas) > 0 {
		for _, v := range datas {
			list, ok := SAPackageMap[v.PackageType]
			if !ok {
				list = make([]cveSa.SAPackageHelper, 0)
			}
			var helper cveSa.SAPackageHelper
			helper.PackageName = v.PackageName
			helper.Child = make([]cveSa.SAPackageHelper, 0)
			list = append(list, helper)
			SAPackageMap[v.PackageType] = list
		}
	}
	for k, v := range SAPackageMap {
		var helper cveSa.SAPackageHelper
		helper.ProductName = k
		helper.Child = v
		SAPackages = append(SAPackages, helper)
	}
	return SAPackages
}

func reReference(datas []models.CveSecurityNoticeReference) []models.RCveSecurityNoticeReference {
	var list = make([]models.RCveSecurityNoticeReference, 0, len(datas))
	for _, v := range datas {
		list = append(list, models.RCveSecurityNoticeReference{CveSecurityNoticeReference: v, Updateime: v.Updateime.Format(_const.Format)})
	}
	return list
}
