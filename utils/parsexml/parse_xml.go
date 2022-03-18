package parsexml

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"cve-sa-backend/iniconf"
	"cve-sa-backend/models"
	"cve-sa-backend/utils"
	_const "cve-sa-backend/utils/const"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"go.uber.org/zap"
)

func GetCVEDatabase(sa, affectedComponent string, child utils.Vulnerability, updateTime time.Time) (cveSa.DatabaseData, error) {

	cve := SetVulnerability(sa, affectedComponent, child, updateTime)

	//bf := bytes.NewBuffer([]byte{})
	//jsonEncoder := json.NewEncoder(bf)
	//jsonEncoder.SetEscapeHTML(false)
	//jsonEncoder.Encode(cve)
	//fmt.Println(bf.String())

	bean, err := GetParserBean(cve.CveId, cve.PackageName, updateTime)
	if err != nil {
		return cve, err
	}

	//fmt.Println(bean.CveParser.Updateime)
	//bf := bytes.NewBuffer([]byte{})
	//jsonEncoder := json.NewEncoder(bf)
	//jsonEncoder.SetEscapeHTML(false)
	//jsonEncoder.Encode(bean)
	//fmt.Println(bf.String())
	cve.CveDatabase.NationalCyberAwarenessSystem = bean.CveParser.Cvss
	if cve.CveDatabase.Type == "" {
		if bean.CveParser.Score != "" {
			score, ok := utils.InterceptString(bean.CveParser.Score, " ", "")
			if ok {
				score = strings.TrimSpace(score)
				score = strings.ToUpper(score[:1]) + strings.ToLower(score[1:])
				cve.CveDatabase.Type = score
			}
		}
	}

	if bean.CveParser.Score != "" {
		score, ok := utils.InterceptString(bean.CveParser.Score, " ", "")
		if ok {
			cve.CveDatabase.CvsssCorenvd = strings.TrimSpace(score)
		}

	}

	if bean.CveParser.Vector != "" {
		vector := utils.GetVectorArr(bean.CveParser.Vector)
		cve.CveDatabase.AttackVectornvd = vector.AV
		cve.CveDatabase.AttackComplexitynvd = vector.AC
		cve.CveDatabase.PrivilegesRequirednvd = vector.PR
		cve.CveDatabase.UserInteractionnvd = vector.UI
		cve.CveDatabase.Scopenvd = vector.S
		cve.CveDatabase.Confidentialitynvd = vector.C
		cve.CveDatabase.Integritynvd = vector.I
		cve.CveDatabase.Availabilitynvd = vector.A
	}
	cve.ParserBean = &bean

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(cve)
	fmt.Println(bf.String())
	return cve, nil
}

type Vulnerability struct {
	Text    string `xml:",innerxml"`
	Ordinal string `xml:"Ordinal,attr"`
	Xmlns   string `xml:"xmlns,attr"`
}

func SetVulnerability(sa, affectedComponent string, child utils.Vulnerability, updateTime time.Time) cveSa.DatabaseData {
	cve := cveSa.DatabaseData{}
	cve.RCveDatabase.CveDatabase.AffectedProduct = sa

	var statusType string
	var productList []string
	if len(child.Notes.Note) > 0 {
		cve.RCveDatabase.CveDatabase.Summary = child.Notes.Note[0].Note
	}

	cve.RCveDatabase.CveDatabase.AnnouncementTime = utils.GetTextTrim(child.ReleaseDate)

	cve.RCveDatabase.CveDatabase.CveId = utils.GetTextTrim(child.Cve)

	if len(child.ProductStatuses.Status) > 0 {
		for _, v := range child.ProductStatuses.Status {
			statusType = v.Type
			for _, nv := range v.ProductID {
				productList = append(productList, nv)
			}
		}
	}

	if len(child.Threats.Threat) > 0 {
		cve.RCveDatabase.CveDatabase.Type = child.Threats.Threat[0].Description
	}

	if len(child.CVSSScoreSets.ScoreSet) > 0 {
		cve.RCveDatabase.CveDatabase.CvsssCoreoe = child.CVSSScoreSets.ScoreSet[0].BaseScore
		if child.CVSSScoreSets.ScoreSet[0].Vector != "" {
			vector := utils.GetVectorArr(child.CVSSScoreSets.ScoreSet[0].Vector)
			cve.RCveDatabase.CveDatabase.AttackVectoroe = vector.AV
			cve.RCveDatabase.CveDatabase.AttackComplexityoe = vector.AC
			cve.RCveDatabase.CveDatabase.PrivilegesRequiredoe = vector.PR
			cve.RCveDatabase.CveDatabase.UserInteractionoe = vector.UI
			cve.RCveDatabase.CveDatabase.Scopeoe = vector.S
			cve.RCveDatabase.CveDatabase.Confidentialityoe = vector.C
			cve.RCveDatabase.CveDatabase.Integrityoe = vector.I
			cve.RCveDatabase.CveDatabase.Availabilityoe = vector.A
		}
	}

	for _, v := range child.Remediations.Remediation {
		if strings.ToLower(statusType) == strings.ToLower("Unaffected") {
			if v.Description != "" {
				affectedComponent = utils.GetTextTrim(v.Description)
			}
		} else {
			break
		}
	}

	var packageList []models.RCveProductPackage
	for _, v := range productList {
		packageObj := models.RCveProductPackage{}
		packageObj.CveProductPackage.CveId = cve.CveId
		packageObj.CveProductPackage.PackageName = affectedComponent
		packageObj.CveProductPackage.ProductName = v
		packageObj.CveProductPackage.Status = statusType
		packageObj.CveProductPackage.Updateime = updateTime
		packageList = append(packageList, packageObj)
	}

	cve.RCveDatabase.CveDatabase.PackageName = affectedComponent
	cve.PackageList = packageList
	cve.RCveDatabase.CveDatabase.Updateime = updateTime

	vu := Vulnerability{}
	vu.Text = child.Text
	vu.Ordinal = child.Ordinal
	vu.Xmlns = child.Xmlns

	by, err := xml.Marshal(&vu)
	if err != nil {
		iniconf.Log.Error("xml.Marshal error:", zap.String("error", err.Error()))
	}

	sb := new(strings.Builder)
	sb.WriteString(_const.TOP + "\n")
	sb.WriteString("<" + _const.CVRFDOC + " xmlns=" + _const.XMLNS + " xmlns:cvrf=" + _const.XmlnsCvrf + ">" + "\n")
	sb.WriteString(string(by) + "\n")
	sb.WriteString("</" + _const.CVRFDOC + ">")

	cvrf := models.RCveCvrf{}
	cvrf.CveCvrf.Cvrf = sb.String()
	cvrf.CveCvrf.CveId = cve.CveId
	cvrf.CveCvrf.PackageName = cve.PackageName
	cvrf.CveCvrf.Updateime = updateTime
	cve.Cvrf = &cvrf

	return cve
}
