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

	"go.uber.org/zap"

	"gorm.io/gorm"
)

type UploadHandle struct {
}

func (u *UploadHandle) DeleteCVE(cveId, packageName string) (string, error) {
	tx := DB.Begin()
	if packageName == "" {
		//Delete the cve information corresponding to the specified cveId
		delCves, rowsAffected, err := dao.DefaultCveDatabase.GetCveDatabaseByCveIdList(&models.CveDatabase{
			CveId: cveId,
		}, tx)
		if err != nil || rowsAffected == 0 {
			tx.Rollback()
			return "", errors.New("CVE error")
		}
		for _, delCve := range delCves {
			err = u.DeleteOneCVE(delCve, tx)
			if err != nil {
				tx.Rollback()
				return "", err
			}
		}
		tx.Commit()
		return fmt.Sprintf("Delete CVE %d record.", len(delCves)), nil
	} else {
		delCve, err := dao.DefaultCveDatabase.GetOneDatabase(&models.CveDatabase{
			CveId:       cveId,
			PackageName: packageName,
		}, tx)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return "", err
		}
		err = u.DeleteOneCVE(*delCve, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		tx.Commit()
		return "Delete CVE 1 record.", nil
	}
}

func (u *UploadHandle) DeleteOneCVE(delCve models.CveDatabase, tx *gorm.DB) error {
	cveCvrf, err := dao.DefaultCvrf.GetOneCvrf(&models.CveCvrf{
		CveId:       delCve.CveId,
		PackageName: delCve.PackageName,
	}, tx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		err = dao.DefaultCvrf.DeleteCvrf(cveCvrf.Id, tx)
		if err != nil {
			return err
		}
	}

	packageList, rowsAffected, err := dao.DefaultCveProductPackage.GetProductPackageList(&models.CveProductPackage{
		CveId:       delCve.CveId,
		PackageName: delCve.PackageName,
	}, tx)
	if err != nil {
		return err
	}
	if rowsAffected > 0 {
		for _, v := range packageList {
			err = dao.DefaultCveProductPackage.DeleteProductPackage(v.Id, tx)
			if err != nil {
				return err
			}
		}
	}

	parserBean, err := dao.DefaultCveParser.GetOneParser(&models.CveParser{
		Cve:         delCve.CveId,
		PackageName: delCve.PackageName,
	}, tx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		err = dao.DefaultCveParser.DeleteParser(parserBean.Id, tx)
		if err != nil {
			return err
		}
	}

	err = dao.DefaultCveDatabase.DeleteCveDatabase(delCve.Id, tx)
	if err != nil {
		return err
	}
	return nil
}

func (u *UploadHandle) DeleteSA(saNo string) error {
	tx := DB.Begin()
	delSecurity, err := dao.DefaultSecurityNotice.GetSecurityNotice(&models.CveSecurityNotice{
		SecurityNoticeNo: saNo,
	}, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	delPackageList, err := dao.DefaultSecurityNoticePackage.GetPackageList(&models.CveSecurityNoticePackage{
		SecurityNoticeNo: delSecurity.SecurityNoticeNo,
	}, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(delPackageList) > 0 {
		err = dao.DefaultSecurityNoticePackage.DeletePackages(delPackageList, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	delReferenceList, err := dao.DefaultSecurityNoticeReference.GetReferenceList(&models.CveSecurityNoticeReference{
		SecurityNoticeNo: delSecurity.SecurityNoticeNo,
	}, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(delPackageList) > 0 {
		err = dao.DefaultSecurityNoticeReference.DeleteReferences(delReferenceList, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	cveCvrf, err := dao.DefaultCvrf.GetOneCvrf(&models.CveCvrf{
		SecurityNoticeNo: delSecurity.SecurityNoticeNo,
	}, tx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return err
	}
	if err == nil {
		err = dao.DefaultCvrf.DeleteCvrf(cveCvrf.Id, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = dao.DefaultSecurityNotice.DeleteSecurityNotice(delSecurity.Id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (u *UploadHandle) GetHttpParserBeanListByCve(cve, packageName string) (*models.CveParser, error) {
	return dao.DefaultCveParser.GetOneParserWithDB(&models.CveParser{
		Cve:         cve,
		PackageName: packageName,
	})
}

func (u *UploadHandle) SyncCve(cveFileName string) (string, error) {
	fileByte, err := utils.GetCvrfFile(cveFileName)
	if err != nil {
		Log.Error("utils.GetCvrfFile error:", zap.Error(err))
		return "SyncCve failed. An exception occurred." + err.Error(), err
	}

	Element := utils.FixedCveXml{}
	err = xml.Unmarshal(fileByte, &Element)
	if err != nil {
		Log.Error("xml.Unmarshal error:", zap.Error(err))
		return "SyncCve failed. An exception occurred." + err.Error(), err
	}

	var list []cveSa.DatabaseData
	updateTime := time.Now()
	for _, v := range Element.Vulnerability {
		cve, err := parsexml.GetCVEDatabase("", "", v, updateTime)
		if err != nil {
			Log.Error("parsexml.GetCVEDatabase error:", zap.Error(err))
			return "SyncCve failed. An exception occurred." + err.Error(), err
		}
		list = append(list, cve)
	}
	tx := DB.Begin()
	err = u.SaveAndDeleteCveList(list, tx)
	if err != nil {
		Log.Error("SaveAndDeleteCveList error:", zap.Error(err))
		tx.Rollback()
		return "SyncCve failed. An exception occurred." + err.Error(), err
	}
	tx.Commit()
	return "CVE sync successfully" + "<br/>", nil
}

func (u *UploadHandle) SaveAndDeleteCveList(list []cveSa.DatabaseData, tx *gorm.DB) error {
	for _, v := range list {
		err := dao.DefaultCvrf.DeleteByCveIdAndPackageName(v.Cvrf.CveId, v.Cvrf.PackageName, tx)
		if err != nil {
			return err
		}
		err = dao.DefaultCvrf.InsertCvrf(&v.Cvrf.CveCvrf, tx)
		if err != nil {
			return err
		}
		status := 0
		for _, pk := range v.PackageList {

			err = dao.DefaultCveProductPackage.DeleteByCveIdAndPackageNameAndProductName(pk.CveId, pk.PackageName, pk.ProductName, tx)
			if err != nil {
				return err
			}
			err = dao.DefaultCveProductPackage.InsertProductPackage(&pk.CveProductPackage, tx)
			if err != nil {
				return err
			}

			if strings.EqualFold(pk.CveProductPackage.Status, "Fixed") {
				status = 1
			}
			if status != 1 {
				if strings.EqualFold(pk.CveProductPackage.Status, "Unaffected") {
					status = 2
				}
			}
		}

		if status == 1 {
			v.RCveDatabase.CveDatabase.Status = "Fixed"
		}
		if status == 2 {
			v.RCveDatabase.CveDatabase.Status = "Unaffected"
		}

		err = dao.DefaultCveParser.DeleteByCveAndPackageName(v.ParserBean.CveParser.Cve, v.ParserBean.CveParser.PackageName, tx)
		if err != nil {
			return err
		}
		err = dao.DefaultCveParser.InsertParser(&v.ParserBean.CveParser, tx)
		if err != nil {
			return err
		}

		err = dao.DefaultCveDatabase.DeleteCveDatabaseByCveIdAndPackageName(v.RCveDatabase.CveDatabase.CveId, v.RCveDatabase.CveDatabase.PackageName, tx)
		if err != nil {
			return err
		}
		err = dao.DefaultCveDatabase.InsertCveDatabase(&v.RCveDatabase.CveDatabase, tx)
		if err != nil {
			return err
		}

	}

	return nil
}

func (u *UploadHandle) SyncHardwareCompatibility() (string, error) {
	var listZh, listEn []*cveSa.HardwareCompatibility
	var err error
	listZh, err = u.parserHardwareCompatibility("zh")
	if err != nil {
		return fmt.Sprint("获取数据失败:"), err
	}
	tx := DB.Begin()
	if listZh != nil && len(listZh) > 0 {
		err = dao.DefaultCompatibilityHardware.DeleteHardwareForLang("zh", tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		u.assemblyHardware(listZh, "zh")
		err = u.saveList(listZh, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	listEn, err = u.parserHardwareCompatibility("en")
	if err != nil {
		return fmt.Sprint("获取数据失败:"), err
	}

	if listEn != nil && len(listEn) > 0 {
		err = dao.DefaultCompatibilityHardware.DeleteHardwareForLang("en", tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		u.assemblyHardware(listEn, "en")
		err = u.saveList(listEn, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}
	tx.Commit()
	return "success", nil
}

// parserHardwareCompatibility Request zh or en json data, and json unmarshal cveSa.HardwareCompatibility
func (u *UploadHandle) parserHardwareCompatibility(lang string) ([]*cveSa.HardwareCompatibility, error) {
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
	var datas []*cveSa.HardwareCompatibility
	err = json.Unmarshal(bytes, &datas)
	if err != nil {
		SLog.Error("json unmarshal failed,err :", err)
		return nil, err
	}
	return datas, nil
}

// assemblyHardware assembly oe_compatibility_hardware insert datas
func (u *UploadHandle) assemblyHardware(datas []*cveSa.HardwareCompatibility, lang string) {
	timeStr := utils.GetCurTime()
	for k := range datas {
		data := datas[k]
		if data.CompatibilityConfiguration == nil {
			continue
		}
		data.OeCompatibilityHardware.ProductInformation = data.CompatibilityConfiguration.ProductInformation
		data.OeCompatibilityHardware.CertificationTime = data.CompatibilityConfiguration.CertificationTime
		data.OeCompatibilityHardware.CommitID = data.CompatibilityConfiguration.CommitID
		data.OeCompatibilityHardware.MotherBoardRevision = data.CompatibilityConfiguration.MotherBoardRevision
		data.OeCompatibilityHardware.BiosUefi = data.CompatibilityConfiguration.BiosUefi
		data.OeCompatibilityHardware.Cpu = data.CompatibilityConfiguration.Cpu
		data.OeCompatibilityHardware.Ram = data.CompatibilityConfiguration.Ram
		data.OeCompatibilityHardware.PortsBusTypes = data.CompatibilityConfiguration.PortsBusTypes
		data.OeCompatibilityHardware.VideoAdapter = data.CompatibilityConfiguration.VideoAdapter
		data.OeCompatibilityHardware.HostBusAdapter = data.CompatibilityConfiguration.HostBusAdapter
		data.OeCompatibilityHardware.HardDiskDrive = data.CompatibilityConfiguration.HardDiskDrive
		data.OeCompatibilityHardware.Updateime = timeStr
		data.OeCompatibilityHardware.Lang = lang
	}
}
func (u *UploadHandle) saveList(datas []*cveSa.HardwareCompatibility, tx *gorm.DB) error {
	for k := range datas {
		hardware := datas[k]
		list := hardware.OeCompatibilityHardware
		err := dao.DefaultCompatibilityHardware.CreateHardware(&list, tx)
		if err != nil {
			return err
		}

		for kk := range hardware.BoardCards {
			adapter := hardware.BoardCards[kk]
			adapter.HardwareId = list.Id
			adapter.Lang = list.Lang
			adapter.Updateime = list.Updateime
			err = dao.DefaultCompatibilityHardwareAdapter.CreateAdapter(adapter, tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *UploadHandle) SyncDriverCompatibility() (string, error) {
	var listZh, listEn []models.OeCompatibilityDriverResponse
	var err error
	listZh, err = u.parserOEDriverCompatibility("zh")
	if err != nil {
		return fmt.Sprint("获取数据失败:"), err
	}
	tx := DB.Begin()
	if listZh != nil && len(listZh) > 0 {
		err = dao.DefaultCompatibilityDriver.DeleteDriverForLang("zh", tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		drivers := u.assemblyDriver(listZh, "zh")
		err = dao.DefaultCompatibilityDriver.CreateDriver(drivers, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	listEn, err = u.parserOEDriverCompatibility("en")
	if err != nil {
		return fmt.Sprint("获取数据失败:"), err
	}
	if listEn != nil && len(listEn) > 0 {
		err = dao.DefaultCompatibilityDriver.DeleteDriverForLang("en", tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		drivers := u.assemblyDriver(listEn, "en")
		err = dao.DefaultCompatibilityDriver.CreateDriver(drivers, tx)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}
	tx.Commit()
	return "success", nil
}

func (u *UploadHandle) parserOEDriverCompatibility(lang string) ([]models.OeCompatibilityDriverResponse, error) {
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

func (u *UploadHandle) assemblyDriver(datas []models.OeCompatibilityDriverResponse, lang string) []models.OeCompatibilityDriver {
	var list = make([]models.OeCompatibilityDriver, 0, len(datas))
	timeStr := utils.GetCurTime()
	for _, v := range datas {
		list = append(list, u.joinDriver(v, timeStr, lang))
	}
	return list
}

func (u *UploadHandle) joinDriver(data models.OeCompatibilityDriverResponse, timeStr, lang string) models.OeCompatibilityDriver {
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

func (u *UploadHandle) TransferData(cve string) (string, error) {
	tx := DB.Begin()
	switch cve {
	case "SAreference":
		list, err := dao.DefaultSecurityNotice.FindAllSecurityNotice()
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
			err = dao.DefaultSecurityNoticeReference.CreateReference(insert, tx)
			if err != nil {
				tx.Rollback()
				return "", err
			}
			tx.Commit()
		}
	case "CVEstatus":
		list, err := dao.DefaultCveDatabase.FindAllCveDatabase()
		if err != nil {
			return "", err
		}
		for _, v := range list {
			packageList, total, err := dao.DefaultCveProductPackage.GetProductPackageList(&models.CveProductPackage{
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
			err = dao.DefaultCveDatabase.UpdateCve(v, tx)
			if err != nil {
				tx.Rollback()
				return "", err
			}
		}
		tx.Commit()
	case "SApackage":
		list, err := dao.DefaultSecurityNotice.FindAllSecurityNotice()
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
					snr.PackageType = u.getPackageType(s)
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
		err = dao.DefaultSecurityNoticePackage.CreatePackage(insert, tx)
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
func (u *UploadHandle) getPackageType(s string) string {
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

func (u *UploadHandle) SyncSA(saFileName string) (string, error) {
	security, err := u.ParserSA(saFileName)
	if err != nil {
		SLog.Error("syncSA ", err)
		return "SyncSA failed. parser exception occurred. error message:" + err.Error(), err
	}
	var delSecurity *models.CveSecurityNotice
	tx := DB.Begin()
	delSecurity, err = dao.DefaultSecurityNotice.GetOneSecurity(&models.CveSecurityNotice{SecurityNoticeNo: security.SecurityNoticeNo}, tx)
	if err != nil {
		SLog.Error("syncSA ", err)
		return rSyncSA(err), err
	}
	if delSecurity != nil {
		err = dao.DefaultSecurityNotice.DeleteSecurity(delSecurity.Id, tx)
		if err != nil {
			tx.Rollback()
			return rSyncSA(err), err
		}
	}
	err = dao.DefaultSecurityNotice.CreateSecurity(security.RCveSecurityNotice.CveSecurityNotice, tx)
	if err != nil {
		SLog.Error("syncSA ", err)
		tx.Rollback()
		return rSyncSA(err), err
	}

	delPackageList, err := dao.DefaultSecurityNoticePackage.GetPackageList(&models.CveSecurityNoticePackage{SecurityNoticeNo: security.SecurityNoticeNo}, tx)
	if err != nil {
		SLog.Error("syncSA ", err)
		tx.Rollback()
		return rSyncSA(err), err
	}
	if len(delPackageList) > 0 {
		err = dao.DefaultSecurityNoticePackage.DeletePackages(delPackageList, tx)
		if err != nil {
			tx.Rollback()
			return rSyncSA(err), err
		}
	}

	if security.PackageList != nil && len(security.PackageList) > 0 {
		var insert = make([]models.CveSecurityNoticePackage, 0, len(security.PackageList))
		var srcInsert = make([]models.CveSecurityNoticePackage, 0, len(security.PackageList))
		var aarIsert = make([]models.CveSecurityNoticePackage, 0, len(security.PackageList))
		var x86insert = make([]models.CveSecurityNoticePackage, 0, len(security.PackageList))

		for _, v := range security.PackageList {
			if v.CveSecurityNoticePackage.PackageType == "src" {
				srcInsert = append(srcInsert, v.CveSecurityNoticePackage)
			} else if v.CveSecurityNoticePackage.PackageType == "noarch" {
				aarIsert = append(aarIsert, v.CveSecurityNoticePackage)
			} else if v.CveSecurityNoticePackage.PackageType == "aarch64" {
				aarIsert = append(aarIsert, v.CveSecurityNoticePackage)
			} else if v.CveSecurityNoticePackage.PackageType == "x86_64" {
				x86insert = append(x86insert, v.CveSecurityNoticePackage)
			}
		}
		insert = append(insert, srcInsert...)
		insert = append(insert, aarIsert...)
		insert = append(insert, x86insert...)
		err = dao.DefaultSecurityNoticePackage.CreatePackage(insert, tx)
		if err != nil {
			tx.Rollback()
			return rSyncSA(err), err
		}
	}

	delReferenceList, err := dao.DefaultSecurityNoticeReference.GetReferenceList(&models.CveSecurityNoticeReference{SecurityNoticeNo: security.SecurityNoticeNo}, tx)
	if err != nil {
		tx.Rollback()
		return rSyncSA(err), err
	}
	if len(delReferenceList) > 0 {
		err = dao.DefaultSecurityNoticeReference.DeleteReferences(delReferenceList, tx)
		if err != nil {
			tx.Rollback()
			return rSyncSA(err), err
		}
	}

	if security.ReferenceList != nil && len(security.ReferenceList) > 0 {
		var i = make([]models.CveSecurityNoticeReference, 0, len(security.ReferenceList))
		for _, v := range security.ReferenceList {
			i = append(i, v.CveSecurityNoticeReference)
		}
		err = dao.DefaultSecurityNoticeReference.CreateReference(i, tx)
		if err != nil {
			tx.Rollback()
			return rSyncSA(err), err
		}
	}

	delCvrf, err := dao.DefaultCvrf.GetOneCvrf(&models.CveCvrf{SecurityNoticeNo: security.SecurityNoticeNo}, tx)
	if err != nil && !utils.ErrorNotFound(err) {
		SLog.Error(err)
		tx.Rollback()
		return rSyncSA(err), err
	}
	if delCvrf.Id > 0 {
		err = dao.DefaultCvrf.DeleteCvrf(delCvrf.Id, tx)
		if err != nil {
			tx.Rollback()
			return rSyncSA(err), err
		}
	}

	if security.Cvrf != nil {
		err = dao.DefaultCvrf.InsertCvrf(&security.Cvrf.CveCvrf, tx)
		if err != nil {
			tx.Rollback()
			return rSyncSA(err), err
		}
	}
	if len(security.CveList) > 0 {
		err = u.SaveAndDeleteCveList(security.CveList, tx)
		if err != nil {
			tx.Rollback()
			return rSyncSA(err), err
		}
	}
	tx.Commit()
	return "SA sync successfully" + _const.BR, nil
}

func rSyncSA(err error) string {
	return "SyncSA failed. database exception occurred. error message:" + err.Error()
}
func (u *UploadHandle) ParserSA(url string) (*cveSa.SecurityNoticeData, error) {
	var security = new(cveSa.SecurityNoticeData)
	var cveList = make([]cveSa.DatabaseData, 0)
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

	err = u.setRevisionHistory(security, element.DocumentTracking)
	if err != nil {
		return nil, err
	}
	u.setSecurity(security, element.DocumentNotes)
	u.setReference(security, element.DocumentReferences, now)
	u.setProduct(security, element.ProductTree, now)
	for _, v := range element.Vulnerability {
		cve, err := parsexml.GetCVEDatabase(security.SecurityNoticeNo, security.AffectedComponent, v, now)
		if err != nil {
			SLog.Error("GetCVEDatabase ", err)
		}
		cveList = append(cveList, cve)
	}
	security.RCveSecurityNotice.CveSecurityNotice.Updateime = now
	security.CveList = cveList

	cvrf := new(models.RCveCvrf)
	cvrf.Cvrf = string(bytes)
	cvrf.SecurityNoticeNo = security.SecurityNoticeNo
	cvrf.CveCvrf.Updateime = now
	security.Cvrf = cvrf
	return security, nil
}

func (u *UploadHandle) setRevisionHistory(security *cveSa.SecurityNoticeData, child utils.DocumentTracking) error {
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

func (u *UploadHandle) setSecurity(security *cveSa.SecurityNoticeData, child utils.DocumentNotes) {
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

func (u *UploadHandle) setReference(security *cveSa.SecurityNoticeData, child utils.DocumentReferences, time time.Time) {
	var referenceList = make([]models.RCveSecurityNoticeReference, 0)
	var cveId = ""
	for _, v := range child.Reference {
		ty := v.Type
		for _, vv := range v.URL {
			snr := models.RCveSecurityNoticeReference{}
			snr.SecurityNoticeNo = security.SecurityNoticeNo
			snr.Type = ty
			snr.Url = vv
			snr.CveSecurityNoticeReference.Updateime = time
			referenceList = append(referenceList, snr)

			if ty == _const.OpenEulerCVE {
				start := strings.Index(vv, "CVE-")
				if start > -1 && len(v.URL) > 0 {
					cveId += vv[start:] + ";"
				}
			}
		}
	}
	security.CveId = cveId
	security.ReferenceList = referenceList
}

func (u *UploadHandle) setProduct(security *cveSa.SecurityNoticeData, child utils.ProductTree, time time.Time) {
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
				snp.ProductName = u.getOrDefault(cpeProductMap, full.CPE, "")
				snp.CveSecurityNoticePackage.Updateime = time
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

func (u *UploadHandle) getOrDefault(m map[string]string, key, defaultV string) string {
	if value, ok := m[key]; ok {
		return value
	}
	return defaultV
}

func rSyncAll() string {
	return "SyncAll failed. An exception occurred."
}
func (u *UploadHandle) SyncAll() (string, error) {
	var result string
	var err error
	tx := DB.Begin()
	SLog.Info("=====>>>>>start sync update_unaffect.txt")
	unaffect, err := utils.GetCvrfFile(_const.UpdateUnaddect)
	if err != nil {
		return rSyncAll(), err
	}
	cvrfFileUnaffect := string(unaffect)

	fixed, err := utils.GetCvrfFile(_const.UpdateFixed)
	if err != nil {
		return rSyncAll(), err
	}
	cvrfFileFixed := string(fixed)

	cveSync, err := dao.DefaultCvrfSync.GetSyncOne(&models.CveCvrfSync{Type: _const.UpdateUnaddect}, tx)
	if err != nil {
		return rSyncAll(), err
	}
	saSync, err := dao.DefaultCvrfSync.GetSyncOne(&models.CveCvrfSync{Type: _const.UpdateFixed}, tx)
	if err != nil {
		return rSyncAll(), err
	}

	if cveSync == nil || cveSync.CvrfFile != cvrfFileUnaffect {
		if cveSync == nil {
			cveSync = new(models.CveCvrfSync)
		}
		SLog.Info("=====>>>>>sync:" + cvrfFileUnaffect)
		if cvrfFileUnaffect != "" && strings.HasSuffix(cvrfFileUnaffect, ".xml") {
			result, err = u.SyncCve(cvrfFileUnaffect)
			if err != nil {
				return rSyncAll(), err
			}
		}
		cveSync.CvrfFile = cvrfFileUnaffect
		cveSync.Type = _const.UpdateUnaddect
		cveSync.Updateime = time.Now()
		if strings.Index(result, "successfully") > -1 {
			err = dao.DefaultCvrfSync.CreateSyncOne(cveSync, tx)
			if err != nil {
				tx.Rollback()
				return rSyncAll(), err
			}
		}
	}
	SLog.Info("=====>>>>>end sync update_unaffect.txt")
	SLog.Info("=====>>>>>start sync update_fixed.txt")

	replaceStr := strings.Replace(cvrfFileFixed, "\n", ",", -1)
	SLog.Info("=====>>>>>sync:" + replaceStr)
	if saSync == nil || saSync.CvrfFile != replaceStr {
		if saSync == nil {
			saSync = new(models.CveCvrfSync)
		}

		if cvrfFileFixed != "" {
			arr := strings.Split(cvrfFileFixed, "\n")
			for _, v := range arr {
				if strings.HasSuffix(v, ".xml") {
					result, _ = u.SyncSA(v)
					if strings.Index(result, "successfully") == -1 {
						SLog.Error("===>>>update_fixed.txt failed ,message:" + result)
						break
					}
				}
			}
		}
		saSync.CvrfFile = replaceStr
		saSync.Type = _const.UpdateFixed
		saSync.Updateime = time.Now()
		if strings.Index(result, "successfully") > -1 {
			err = dao.DefaultCvrfSync.CreateSyncOne(saSync, tx)
			if err != nil {
				tx.Rollback()
				return rSyncAll(), err
			}
		}
	}
	SLog.Info("=====>>>>>end sync update_fixed.txt")
	tx.Commit()
	return result, nil
}

func (u *UploadHandle) SyncOsv() (string, error) {
	var tools, platform []byte
	var ok bool
	t := time.Now()
	osvList, err := u.parserOsv()
	if err != nil {
		return fmt.Sprint("获取数据失败:"), err
	}
	tx := DB.Begin()

	for k := range osvList {
		v := osvList[k]
		if len(v.PlatformResult) == 0 && len(v.ToolsResult) == 0 {
			err = dao.DefaultCompatibilityOsv.DeleteOsv(v.OsVersion, tx)
			if err != nil {
				tx.Rollback()
				return "", err
			}
			continue
		}
		tools, err = json.Marshal(v.ToolsResult)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		platform, err = json.Marshal(v.PlatformResult)
		if err != nil {
			tx.Rollback()
			return "", err
		}
		osv := models.OeCompatibilityOsv{
			Architecture:         v.Arch,
			OsVersion:            v.OsVersion,
			OsvName:              v.OsvName,
			Date:                 v.Date,
			OsDownloadLink:       v.OsDownloadLink,
			Type:                 v.Type,
			Details:              v.Details,
			FriendlyLink:         v.FriendlyLink,
			TotalResult:          v.TotalResult,
			CheckSum:             v.CheckSum,
			BaseOpeneulerVersion: v.BaseOpeneulerVersion,
			ToolsResult:          string(tools),
			PlatformResult:       string(platform),
			Updateime:            t,
		}
		if ok, err = dao.DefaultCompatibilityOsv.ExistsOsv(osv.OsVersion, tx); err == nil && ok {
			err = dao.DefaultCompatibilityOsv.UpdateOsv(osv, tx)
			if err != nil {
				tx.Rollback()
				return "", err
			}
		} else if err == nil {
			err = dao.DefaultCompatibilityOsv.CreateOsv(osv, tx)
			if err != nil {
				tx.Rollback()
				return "", err
			}
		} else {
			tx.Rollback()
			return "", err
		}
	}
	tx.Commit()
	return "success", nil
}

func (u *UploadHandle) parserOsv() ([]cveSa.Osv, error) {
	bytes, err := utils.HTTPGet(_const.ParserOsvJsonFile)
	if err != nil {
		return nil, err
	}

	var data []cveSa.Osv

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		SLog.Error("Osv json unmarshal failed, err :", err)
		return nil, err
	}
	return data, nil
}
