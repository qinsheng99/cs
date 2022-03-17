package manage

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"

	"cve-sa-backend/dao"
	. "cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	_const "cve-sa-backend/utils/const"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/parsexml"

	"gorm.io/gorm"
)

func DeleteCVE(cveId, packageName string) (string, error) {
	tx := DB.Begin()
	if packageName == "" {
		//Delete the cve information corresponding to the specified cveId
		delCves, rowsAffected, err := dao.GetCveDatabaseByCveIdList(&models.CveDatabase{
			CveId: cveId,
		}, tx)
		if err != nil || rowsAffected == 0 {
			return "", errors.New("CVE error")
		}
		for _, delCve := range delCves {
			err = DeleteOneCVE(delCve, tx)
			if err != nil {
				tx.Rollback()
				return "", err
			}
		}
		tx.Commit()
		return fmt.Sprintf("Delete CVE %d record.", len(delCves)), nil
	} else {
		delCve, err := dao.GetOneDatabase(&models.CveDatabase{
			CveId:       cveId,
			PackageName: packageName,
		}, tx)
		if err == nil {
			err = DeleteOneCVE(*delCve, tx)
			if err != nil {
				tx.Rollback()
				return "", err
			}
			tx.Commit()
			return "Delete CVE 1 record.", nil
		} else {
			tx.Rollback()
			return "", err
		}
	}
}

func DeleteOneCVE(delCve models.CveDatabase, tx *gorm.DB) error {
	cveCvrf, err := dao.GetOneCvrf(&models.CveCvrf{
		CveId:       delCve.CveId,
		PackageName: delCve.PackageName,
	}, tx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		err = dao.DeleteCvrf(cveCvrf.Id, tx)
		if err != nil {
			return err
		}
	}

	packageList, rowsAffected, err := dao.GetProductPackageList(&models.CveProductPackage{
		CveId:       delCve.CveId,
		PackageName: delCve.PackageName,
	}, tx)
	if err != nil {
		return err
	}
	if rowsAffected > 0 {
		for _, v := range packageList {
			err = dao.DeleteProductPackage(v.Id, tx)
			if err != nil {
				return err
			}
		}
	}

	parserBean, err := dao.GetOneParser(&models.CveParser{
		Cve:         delCve.CveId,
		PackageName: delCve.PackageName,
	}, tx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		err = dao.DeleteParser(parserBean.Id, tx)
		if err != nil {
			return err
		}
	}

	err = dao.DeleteCveDatabase(delCve.Id, tx)
	if err != nil {
		return err
	}
	return nil
}

func SyncCve(cveFileName string) error {
	list := make([]models.CveDatabase, 0, 10)
	fileByte, err := utils.GetCvrfFile(cveFileName)
	if err != nil {
		return err
	}

	Element := utils.FixedCveXml{}
	err = xml.Unmarshal(fileByte, &Element)
	if err != nil {
		return err
	}
	updateTime := time.Now()
	for _, v := range Element.Vulnerability {
		cve, err := parsexml.GetCVEDatabase("", "", v, updateTime)
		if err != nil {
			return err
		}
		fmt.Println(cve)
		list = append(list, models.CveDatabase{})
	}

	return nil
}

func SyncHardwareCompatibility() (string, error) {
	var listZh, listEn []cveSa.HardwareCompatibility
	var err error
	listZh, err = parserHardwareCompatibility("zh")
	if err != nil {
		return fmt.Sprint("获取数据失败:", err.Error()), err
	}
	tx := DB.Begin()
	if listZh != nil && len(listZh) > 0 {
		err = dao.DeleteHardwareForLang("zh", tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		hardwareds := assemblyHardware(listZh, "zh")
		err = dao.CreateHardware(hardwareds, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	listEn, err = parserHardwareCompatibility("en")
	if err != nil {
		return fmt.Sprint("获取数据失败:", err.Error()), err
	}

	if listEn != nil && len(listEn) > 0 {
		err = dao.DeleteHardwareForLang("en", tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		hardwareds := assemblyHardware(listEn, "en")
		err = dao.CreateHardware(hardwareds, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}
	tx.Commit()
	return "success", nil
}

// parserHardwareCompatibility Request zh or en json data, and json unmarshal cveSa.HardwareCompatibility
func parserHardwareCompatibility(lang string) ([]cveSa.HardwareCompatibility, error) {
	var bytes []byte
	var err error
	switch lang {
	case "zh":
		bytes, err = utils.HTTPGet(_const.ParserHardwareFileZh)
		if err != nil {
			SLog.Error("http request failed, err :", err)
			return nil, err
		}
	case "en":
		bytes, err = utils.HTTPGet(_const.ParserHardwareFileEn)
		if err != nil {
			SLog.Error("http request failed, err :", err)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("the input parameter `lang` must be zh or en")
	}
	if strings.Index(string(bytes), "<!DOCTYPE html>") > -1 || strings.Index(string(bytes), "<html>") > -1 {
		return nil, nil
	}
	var datas []cveSa.HardwareCompatibility
	err = json.Unmarshal(bytes, &datas)
	if err != nil {
		SLog.Error("json unmarshal failed,err :", err)
		return nil, err
	}
	return datas, nil
}

// assemblyHardware assembly oe_compatibility_hardware insert datas
func assemblyHardware(datas []cveSa.HardwareCompatibility, lang string) []models.OeCompatibilityHardware {
	var list = make([]models.OeCompatibilityHardware, 0, len(datas))
	timeStr := utils.GetCurTime()
	for _, v := range datas {
		if v.CompatibilityConfiguration == nil {
			list = append(list, v.OeCompatibilityHardware)
			continue
		}
		list = append(list, joinHardware(v, lang, timeStr))
	}
	return list
}

// joinHardware
func joinHardware(data cveSa.HardwareCompatibility, lang, timeStr string) models.OeCompatibilityHardware {
	hardware := data.OeCompatibilityHardware
	hardware.ProductInformation = data.CompatibilityConfiguration.ProductInformation
	hardware.CertificationTime = data.CompatibilityConfiguration.CertificationTime
	hardware.CommitID = data.CompatibilityConfiguration.CommitID
	hardware.MotherBoardRevision = data.CompatibilityConfiguration.MotherBoardRevision
	hardware.BiosUefi = data.CompatibilityConfiguration.BiosUefi
	hardware.Cpu = data.CompatibilityConfiguration.Cpu
	hardware.Ram = data.CompatibilityConfiguration.Ram
	hardware.PortsBusTypes = data.CompatibilityConfiguration.PortsBusTypes
	hardware.VideoAdapter = data.CompatibilityConfiguration.VideoAdapter
	hardware.HostBusAdapter = data.CompatibilityConfiguration.HostBusAdapter
	hardware.HardDiskDrive = data.CompatibilityConfiguration.HardDiskDrive
	hardware.Updateime = timeStr
	hardware.Lang = lang
	return hardware
}

func SyncDriverCompatibility() (string, error) {
	var listZh, listEn []models.OeCompatibilityDriverResponse
	var err error
	listZh, err = parserOEDriverCompatibility("zh")
	if err != nil {
		return fmt.Sprint("获取数据失败:", err.Error()), err
	}
	tx := DB.Begin()
	if listZh != nil && len(listZh) > 0 {
		err = dao.DeleteDriverForLang("zh", tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		drivers := assemblyDriver(listZh, "zh")
		err = dao.CreateDriver(drivers, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	listEn, err = parserOEDriverCompatibility("en")
	if err != nil {
		return fmt.Sprint("获取数据失败:", err.Error()), err
	}
	if listEn != nil && len(listEn) > 0 {
		err = dao.DeleteDriverForLang("en", tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		drivers := assemblyDriver(listEn, "en")
		err = dao.CreateDriver(drivers, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}
	tx.Commit()
	return "success", nil
}

func parserOEDriverCompatibility(lang string) ([]models.OeCompatibilityDriverResponse, error) {
	var bytes []byte
	var err error
	switch lang {
	case "zh":
		bytes, err = utils.HTTPGet(_const.ParserDriverFileZh)
		if err != nil {
			SLog.Error("http request failed, err :", err)
			return nil, err
		}
	case "en":
		bytes, err = utils.HTTPGet(_const.ParserDriverFileEn)
		if err != nil {
			SLog.Error("http request failed, err :", err)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("the input parameter `lang` must be zh or en")
	}
	if strings.Index(string(bytes), "<!DOCTYPE html>") > -1 {
		return nil, nil
	}
	var datas []models.OeCompatibilityDriverResponse
	err = json.Unmarshal(bytes, &datas)
	if err != nil {
		SLog.Error("json unmarshal failed,err :", err)
		return nil, err
	}
	return datas, nil
}

func assemblyDriver(datas []models.OeCompatibilityDriverResponse, lang string) []models.OeCompatibilityDriver {
	var list = make([]models.OeCompatibilityDriver, 0, len(datas))
	timeStr := utils.GetCurTime()
	for _, v := range datas {
		list = append(list, joinDriver(v, timeStr, lang))
	}
	return list
}

func joinDriver(data models.OeCompatibilityDriverResponse, timeStr, lang string) models.OeCompatibilityDriver {
	return models.OeCompatibilityDriver{
		Architecture: data.Architecture,
		BoardModel:   data.BoardModel,
		ChipModel:    utils.InterfaceToString(data.ChipModel),
		ChipVendor:   data.ChipVendor,
		Deviceid:     utils.InterfaceToString(data.Deviceid),
		DownloadLink: data.DownloadLink,
		DriverDate:   data.DriverDate,
		DriverName:   data.DriverName,
		DriverSize:   data.DriverSize,
		Item:         data.Item,
		Os:           data.Os,
		Sha256:       data.Sha256,
		SsID:         utils.InterfaceToString(data.SsID),
		SvID:         utils.InterfaceToString(data.SvID),
		Type:         data.Type,
		Vendorid:     utils.InterfaceToString(data.Vendorid),
		Version:      data.Version,
		Lang:         lang,
		Updateime:    timeStr,
	}
}

func TransferData(cve string) (string, error) {
	tx := DB.Begin()
	switch cve {
	case "SAreference":
		list, err := dao.FindAllSecurityNotice()
		if err != nil {
			return "", err
		}
		var insert = make([]models.CveSecurityNoticeReference, 0)
		for _, v := range list {
			if v.ReferenceDocuments != "" {
				arr := strings.Split(v.ReferenceDocuments, "\n")
				for _, s := range arr {
					insert = append(insert, models.CveSecurityNoticeReference{
						SecurityNoticeNo: v.SecurityNoticeNo,
						Type:             _const.TypeOther,
						Url:              s,
						Updateime:        v.Updateime,
					})
				}
			}
		}
		if len(insert) > 0 {
			err = dao.CreateReference(insert, tx)
			if err != nil {
				tx.Rollback()
				return "", err
			}
			tx.Commit()
		}
	case "CVEstatus":
		list, err := dao.FindAllCveDatabase()
		if err != nil {
			return "", err
		}
		for _, v := range list {
			packageList, total, err := dao.GetProductPackageList(&models.CveProductPackage{
				CveId:       v.CveId,
				PackageName: v.PackageName,
			}, DB)
			if err != nil {
				return "", err
			}
			var status = ""

			if total > 0 {
				for _, pl := range packageList {
					if strings.EqualFold(pl.Status, "Fixed") {
						status = "Fixed"
					} else if strings.EqualFold(pl.Status, "Unaffected") {
						//if status is not Fixed,the status is Unaffected
						if status != "Fixed" {
							status = "Unaffected"
						}
					}
				}
			}
			v.Status = status
			err = dao.UpdateCve(v, tx)
			if err != nil {
				tx.Rollback()
				return "", err
			}
		}
		tx.Commit()
	case "SApackage":
		list, err := dao.FindAllSecurityNotice()
		if err != nil {
			return "", err
		}
		var insert = make([]models.CveSecurityNoticePackage, 0, len(list))
		var srcInsert = make([]models.CveSecurityNoticePackage, 0, len(list))
		var aarIsert = make([]models.CveSecurityNoticePackage, 0, len(list))
		var x86insert = make([]models.CveSecurityNoticePackage, 0, len(list))

		for _, v := range list {
			if v.PackageName != "" {
				arr := strings.Split(v.PackageName, ";")
				for _, s := range arr {
					if s == "" || len(s) < 3 {
						continue
					}
					snr := models.CveSecurityNoticePackage{}
					snr.SecurityNoticeNo = v.SecurityNoticeNo
					snr.PackageType = getPackageType(s)
					snr.ProductName = v.AffectedProduct
					snr.PackageName = utils.TrimStringNR(s)
					snr.Updateime = v.Updateime

					if snr.PackageType == "src" {
						srcInsert = append(srcInsert, snr)
					} else if snr.PackageType == "noarch" {
						aarch64 := snr
						aarch64.PackageType = "aarch64"
						aarch64.PackageName = utils.TrimStringNR(s)
						aarIsert = append(aarIsert, aarch64)

						x86 := snr
						x86.PackageType = "x86_64"
						x86.PackageName = utils.TrimStringNR(s)
						x86insert = append(x86insert, x86)
					} else if snr.PackageType == "aarch64" {
						aarIsert = append(aarIsert, snr)
					} else if snr.PackageType == "x86_64" {
						x86insert = append(x86insert, snr)
					}
				}
				insert = append(insert, srcInsert...)
				insert = append(insert, aarIsert...)
				insert = append(insert, x86insert...)
			}
		}
		err = dao.CreatePackage(insert, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		tx.Commit()
	default:
		return "Parameter error", nil
	}
	return "success", nil
}
func getPackageType(s string) string {
	if strings.Contains(s, ".src.rpm") {
		return "src"
	} else if strings.Contains(s, "aarch64.rpm") {
		return "aarch64"
	} else if strings.Contains(s, "x86_64.rpm") {
		return "x86_64"
	} else if strings.Contains(s, "noarch.rpm") {
		return "noarch"
	} else {
		return ""
	}
}

func SyncSA(saFileName string) string {
	sb := new(strings.Builder)
	sb.WriteString("SA sync successfully")
	sb.WriteString(_const.BR)

	security, err := ParserSA(saFileName)
	if err != nil {
		SLog.Error("syncSA ", err)
		return rSyncSA(err)
	}

	var delSecurity *models.CveSecurityNotice
	tx := DB.Begin()
	delSecurity, err = dao.GetOneSecurity(&models.CveSecurityNotice{SecurityNoticeNo: security.SecurityNoticeNo}, tx)
	if err != nil {
		SLog.Error("syncSA ", err)
		return rSyncSA(err)
	}
	if delSecurity != nil {
		err = dao.DeleteSecurity(delSecurity, tx)
		if err != nil {
			tx.Rollback()
			return rSyncSA(err)
		}
	}
	securityData := security.RCveSecurityNotice.CveSecurityNotice
	securityData.Updateime = utils.StrToTime(security.RCveSecurityNotice.Updateime)
	err = dao.CreateSecurity(securityData, tx)
	if err != nil {
		tx.Rollback()
		return rSyncSA(err)
	}
	return sb.String()
}

func rSyncSA(err error) string {
	return "SyncSA failed. parser exception occurred. error message:" + err.Error()
}
func ParserSA(url string) (*cveSa.SecurityNoticeData, error) {
	var security = new(cveSa.SecurityNoticeData)
	now := time.Now()
	bytes, err := utils.GetCvrfFile(url)
	if err != nil {
		return nil, err
	}

	var element utils.FixedCveXml
	err = xml.Unmarshal(bytes, &element)
	if err != nil {
		return nil, err
	}

	err = setRevisionHistory(security, element.DocumentTracking)
	if err != nil {
		return nil, err
	}
	setSecurity(security, element.DocumentNotes)
	setReference(security, element.DocumentReferences, now)
	setProduct(security, element.ProductTree, now)
	security.RCveSecurityNotice.Updateime = now.Format(_const.Format)
	security.CveList = make([]cveSa.DatabaseData, 0)

	cvrf := new(models.RCveCvrf)
	cvrf.Cvrf = string(bytes)
	cvrf.SecurityNoticeNo = security.SecurityNoticeNo
	cvrf.Updateime = now.Format(_const.Format)
	security.Cvrf = cvrf
	return security, nil
}

func setRevisionHistory(security *cveSa.SecurityNoticeData, child utils.DocumentTracking) error {
	security.SecurityNoticeNo = child.Identification.ID
	var listm = make([]map[string]string, 0)
	for _, v := range child.RevisionHistory.Revision {
		var m = make(map[string]string)
		m["Number"] = v.Number
		m["Date"] = v.Date
		m["Description"] = v.Description
		listm = append(listm, m)
	}
	bys, err := json.Marshal(&listm)
	if err != nil {
		SLog.Error("json marshal field,", err)
		return err
	}
	security.RevisionHistory = string(bys)
	security.AnnouncementTime = child.InitialReleaseDate
	return nil
}

func setSecurity(security *cveSa.SecurityNoticeData, child utils.DocumentNotes) {
	for _, v := range child.Note {
		if v.Title == _const.Synopsis {
			security.Summary = v.Note
		}
		if v.Title == _const.Summary {
			security.Introduction = v.Note
		}
		if v.Title == _const.Description {
			security.Description = v.Note
		}
		if v.Title == _const.Topic {
			security.Subject = v.Note
		}
		if v.Title == _const.Severity {
			security.Type = v.Note
		}
		if v.Title == _const.AffectedComponent {
			security.AffectedComponent = v.Note
		}
	}
}

func setReference(security *cveSa.SecurityNoticeData, child utils.DocumentReferences, time time.Time) {
	var referenceList = make([]models.RCveSecurityNoticeReference, 0)
	var cveId = ""
	for _, v := range child.Reference {
		ty := v.Type
		snr := models.RCveSecurityNoticeReference{}
		snr.SecurityNoticeNo = security.SecurityNoticeNo
		snr.Type = ty
		snr.Url = v.URL
		snr.Updateime = time.Format(_const.Format)
		referenceList = append(referenceList, snr)

		if ty == _const.OpenEulerCVE {
			start := strings.Index(v.URL, "CVE-")
			if start > -1 && len(v.URL) > 0 {
				cveId += v.URL[start:] + ";"
			}
		}
	}
	security.CveId = cveId
	security.ReferenceList = referenceList
}

func setProduct(security *cveSa.SecurityNoticeData, child utils.ProductTree, time time.Time) {
	var packageList = make([]models.RCveSecurityNoticePackage, 0)
	var productName = ""
	var cpeProductMap = make(map[string]string, 8)
	for _, v := range child.Branch {
		name := v.Name
		if name == _const.OpenEuler {
			for _, full := range v.FullProductName {
				cpeProductMap[full.CPE] = full.ProductID
				productName += full.ProductName + ";"
			}
		} else {
			for _, full := range v.FullProductName {
				snp := models.RCveSecurityNoticePackage{}
				snp.SecurityNoticeNo = security.SecurityNoticeNo
				snp.PackageName = full.ProductName
				snp.PackageType = name
				snp.ProductName = getOrDefault(cpeProductMap, full.CPE, "")
				snp.Updateime = time.Format(_const.Format)
				packageList = append(packageList, snp)
			}
		}
	}
	if len(productName) > 1 {
		productName = productName[0 : len(productName)-1]
	}
	security.AffectedProduct = productName
	security.PackageList = packageList
}

func getOrDefault(m map[string]string, key, defaultV string) string {
	if value, ok := m[key]; ok {
		return value
	}
	return defaultV
}
