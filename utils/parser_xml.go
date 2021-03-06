package utils

import (
	"encoding/xml"
)

type (
	UnaffectCveXml struct {
		Cvrfdoc       xml.Name        `xml:"cvrfdoc"`
		Xmlns         string          `xml:"xmlns,attr"`
		Vulnerability []Vulnerability `xml:"Vulnerability"`
	}

	FixedCveXml struct {
		Xmlns              string             `xml:"xmlns,attr"`
		Cvrf               string             `xml:"cvrf,attr"`
		Text               string             `xml:",innerxml"`
		DocumentTitle      DocumentTitle      `xml:"DocumentTitle"`
		DocumentType       string             `xml:"DocumentType"`
		DocumentPublisher  DocumentPublisher  `xml:"DocumentPublisher"`
		DocumentTracking   DocumentTracking   `xml:"DocumentTracking"`
		DocumentNotes      DocumentNotes      `xml:"DocumentNotes"`
		DocumentReferences DocumentReferences `xml:"DocumentReferences"`
		ProductTree        ProductTree        `xml:"ProductTree"`
		Vulnerability      []Vulnerability    `xml:"Vulnerability"`
	}

	DocumentTitle struct {
		XmlLang       string `xml:"lang,attr"`
		DocumentTitle string `xml:",innerxml"`
	}

	Vulnerability struct {
		Text            string      `xml:",innerxml"`
		Ordinal         string      `xml:"Ordinal,attr"`
		Xmlns           string      `xml:"xmlns,attr"`
		Notes           Note        `xml:"Notes"`
		ReleaseDate     string      `xml:"ReleaseDate,omitempty"`
		Cve             string      `xml:"CVE"`
		ProductStatuses Status      `xml:"ProductStatuses"`
		Threats         Threat      `xml:"Threats,omitempty"`
		CVSSScoreSets   ScoreSet    `xml:"CVSSScoreSets"`
		Remediations    Remediation `xml:"Remediations"`
	}

	Note struct {
		Note []Notes
	}

	Status struct {
		Status []ProductID `xml:"Status"`
	}

	ProductID struct {
		Type      string   `xml:"Type,attr"`
		ProductID []string `xml:"ProductID"`
	}

	ScoreSet struct {
		ScoreSet []ScoreSetChild `xml:"ScoreSet"`
	}

	ScoreSetChild struct {
		BaseScore string `xml:"BaseScore"`
		Vector    string `xml:"Vector"`
	}

	Remediation struct {
		Remediation []RemediationChild `xml:"Remediation"`
	}

	RemediationChild struct {
		Type        string `xml:"Type,attr"`
		Description string `xml:"Description"`
		DATE        string `xml:"DATE"`
		ProductID   string `xml:"ProductID,omitempty"`
		URL         string `xml:"URL,omitempty"`
	}

	DocumentPublisher struct {
		Type             string `xml:"Type,attr"`
		ContactDetails   string `xml:"ContactDetails"`
		IssuingAuthority string `xml:"IssuingAuthority"`
	}

	DocumentTracking struct {
		Identification     IdentificationChild `xml:"Identification"`
		Status             string              `xml:"Status"`
		Version            string              `xml:"Version"`
		RevisionHistory    Revision            `xml:"RevisionHistory"`
		InitialReleaseDate string              `xml:"InitialReleaseDate"`
		CurrentReleaseDate string              `xml:"CurrentReleaseDate"`
		Generator          GeneratorChild      `xml:"Generator"`
	}

	Revision struct {
		Revision []RevisionChild `xml:"Revision"`
	}

	RevisionChild struct {
		Number      string `xml:"Number"`
		Date        string `xml:"Date"`
		Description string `xml:"Description"`
	}

	GeneratorChild struct {
		Engine string `xml:"Engine"`
		Date   string `xml:"Date"`
	}

	IdentificationChild struct {
		ID string `xml:"ID"`
	}

	DocumentNotes struct {
		Note []Notes
	}
	Notes struct {
		Title   string `xml:"Title,attr"`
		Type    string `xml:"Type,attr"`
		Ordinal string `xml:"Ordinal,attr"`
		XmlLang string `xml:"lang,attr"`
		Note    string `xml:",innerxml"`
	}

	DocumentReferences struct {
		Reference []ReferenceChild `xml:"Reference"`
	}

	ReferenceChild struct {
		Type string   `xml:"Type,attr"`
		URL  []string `xml:"URL"`
	}

	ProductTree struct {
		Branch []BranchChild `xml:"Branch"`
	}

	BranchChild struct {
		Type            string                 `xml:"Type,attr"`
		Name            string                 `xml:"Name,attr"`
		FullProductName []FullProductNameChild `xml:"FullProductName"`
	}
	FullProductNameChild struct {
		ProductID   string `xml:"ProductID,attr"`
		CPE         string `xml:"CPE,attr"`
		ProductName string `xml:",innerxml"`
	}

	Threat struct {
		Threat []ThreatChild `xml:"Threat"`
	}

	ThreatChild struct {
		Type        string `xml:"Type,attr"`
		Description string `xml:"Description"`
	}
)
