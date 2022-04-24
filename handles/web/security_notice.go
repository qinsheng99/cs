package web

import (
	"strings"

	"cve-sa-backend/dao"
	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	_const "cve-sa-backend/utils/const"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
)

type SecurityHandle struct {
}

func (s *SecurityHandle) FindAllSecurity(req cveSa.RequestData) (*cveSa.ResultData, error) {
	datas, total, err := dao.DefaultSecurityNotice.SecurityFindAll(req)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}
	return s.returnNoticeData(datas, total), nil
}

func (s *SecurityHandle) returnNoticeData(datas []models.CveSecurityNotice, total int64) *cveSa.ResultData {
	securityNoticeData := s.SecurityNoticeData(datas)
	return &cveSa.ResultData{
		SecurityNoticeList:  securityNoticeData,
		CveDatabaseList:     make([]cveSa.DatabaseData, 0),
		ApplicationCompList: make([]models.ROeCompatibilityApplication, 0),
		HardwareCompList:    make([]cveSa.HardwareCompatibility, 0),
		DriverCompList:      make([]models.OeCompatibilityDriver, 0),
		Total:               int(total),
	}
}

func (s *SecurityHandle) GetSecurityNoticePackageByPackageName(pname string) ([]models.RCveSecurityNoticePackage, error) {
	pname = strings.Replace(pname, "\n", "", -1)
	pnames := strings.Split(pname, ",")

	datas, err := dao.DefaultSecurityNoticePackage.ByPackageName(pnames)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}
	return s.packageData(datas), nil
}

func (s *SecurityHandle) packageData(datas []models.CveSecurityNoticePackage) []models.RCveSecurityNoticePackage {
	var list = make([]models.RCveSecurityNoticePackage, 0, len(datas))
	for _, v := range datas {
		list = append(list, models.RCveSecurityNoticePackage{CveSecurityNoticePackage: v, Updateime: v.Updateime.Format(_const.Format)})
	}
	return list
}

func (s *SecurityHandle) NoticeByCVEID(cveId string) ([]cveSa.SecurityNoticeData, error) {
	datas, err := dao.DefaultSecurityNotice.NoticeByCveId(cveId)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}
	return s.SecurityNoticeData(datas), nil
}

func (s *SecurityHandle) SecurityNoticeData(datas []models.CveSecurityNotice) []cveSa.SecurityNoticeData {
	var securityNoticeData = make([]cveSa.SecurityNoticeData, 0, len(datas))

	for _, v := range datas {
		securityNoticeData = append(securityNoticeData, s.SecurityNoticeDataOne(v))
	}
	return securityNoticeData
}

func (s *SecurityHandle) SecurityNoticeDataOne(datas models.CveSecurityNotice) cveSa.SecurityNoticeData {
	return cveSa.SecurityNoticeData{
		RCveSecurityNotice: models.RCveSecurityNotice{
			CveSecurityNotice: datas,
			Updateime:         datas.Updateime.Format(_const.Format),
		},
		PackageHelperList: make([]cveSa.SAPackageHelper, 0),
		PackageList:       make([]models.RCveSecurityNoticePackage, 0),
		ReferenceList:     make([]models.RCveSecurityNoticeReference, 0),
		CveList:           make([]cveSa.DatabaseData, 0),
	}
}

func (s *SecurityHandle) ByCveIdAndAffectedComponent(cveId, affectedComponent string) ([]cveSa.SecurityNoticeData, error) {
	datas, err := dao.DefaultSecurityNotice.NoticeByCveIdComponent(cveId, affectedComponent)
	if err != nil {
		iniconf.SLog.Error(err)
		return nil, err
	}
	return s.SecurityNoticeData(datas), nil
}

func (s *SecurityHandle) NoticeBySecurityNoticeNo(sec string) (*cveSa.SecurityNoticeData, error) {
	var SAPackages = make([]cveSa.SAPackageHelper, 0)
	securityNotice, err := dao.DefaultSecurityNotice.NoticeByNo(sec)
	if err != nil {
		return nil, err
	}
	if securityNotice == nil {
		return nil, nil
	}
	snData := s.SecurityNoticeDataOne(*securityNotice)
	if securityNotice.AffectedProduct != "" {
		products := strings.Split(securityNotice.AffectedProduct, ";")
		for _, v := range products {
			packages, err := dao.DefaultSecurityNoticePackage.NoticePackageByNoProduct(sec, v)
			if err != nil {
				iniconf.SLog.Error(err)
				return &snData, nil
			}
			var SAPackage cveSa.SAPackageHelper
			SAPackage.ProductName = v
			SAPackage.Child = s.getSAPackageHelper(packages)
			SAPackages = append(SAPackages, SAPackage)
		}
	}
	snData.PackageHelperList = SAPackages
	snData.Description = strings.Replace(snData.Description, "\\r\\n", "\r\n", -1)
	snData.Subject = strings.Replace(snData.Subject, "\\r\\n", "\r\n", -1)
	references, err := dao.DefaultSecurityNoticeReference.GetReferenceByNo(sec)
	if err != nil {
		iniconf.SLog.Error(err)
		return &snData, nil
	}
	snData.ReferenceList = s.reReference(references)
	return &snData, nil
}

func (s *SecurityHandle) getSAPackageHelper(datas []models.CveSecurityNoticePackage) []cveSa.SAPackageHelper {
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

func (s *SecurityHandle) reReference(datas []models.CveSecurityNoticeReference) []models.RCveSecurityNoticeReference {
	var list = make([]models.RCveSecurityNoticeReference, 0, len(datas))
	for _, v := range datas {
		list = append(list, models.RCveSecurityNoticeReference{CveSecurityNoticeReference: v, Updateime: v.Updateime.Format(_const.Format)})
	}
	return list
}
