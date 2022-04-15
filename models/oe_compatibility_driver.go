package models

type OeCompatibilityDriver struct {
	Id           int64  `gorm:"column:id" json:"id"`
	Architecture string `gorm:"column:architecture" json:"architecture"`
	BoardModel   string `gorm:"column:board_model" json:"boardModel"`
	ChipModel    string `gorm:"column:chip_model" json:"chipModel"`
	ChipVendor   string `gorm:"column:chip_vendor" json:"chipVendor"`
	Deviceid     string `gorm:"column:deviceid" json:"deviceID"`
	DownloadLink string `gorm:"column:download_link" json:"downloadLink"`
	DriverDate   string `gorm:"column:driver_date" json:"driverDate"`
	DriverName   string `gorm:"column:driver_name" json:"driverName"`
	DriverSize   string `gorm:"column:driver_size" json:"driverSize"`
	Item         string `gorm:"column:item" json:"item"`
	Lang         string `gorm:"column:lang" json:"lang"`
	Os           string `gorm:"column:os" json:"os"`
	Sha256       string `gorm:"column:sha256" json:"sha256"`
	SsID         string `gorm:"column:ssid" json:"ssID"`
	SvID         string `gorm:"column:svid" json:"svID"`
	Type         string `gorm:"column:type" json:"type"`
	Vendorid     string `gorm:"column:vendorid" json:"vendorID"`
	Version      string `gorm:"column:version" json:"version"`
	Updateime    string `gorm:"column:update_time" json:"updateTime"`
}

func (o *OeCompatibilityDriver) TableName() string {
	return "oe_compatibility_driver"
}

type ROeCompatibilityDriver struct {
	OeCompatibilityDriver
	Updateime string `json:"updateTime"`
}

type OeCompatibilityDriverResponse struct {
	Id           int64       `json:"id"`
	Architecture string      `json:"architecture"`
	BoardModel   string      `json:"boardModel"`
	ChipModel    interface{} `json:"chipModel"`
	ChipVendor   string      `json:"chipVendor"`
	Deviceid     interface{} `json:"deviceID"`
	DownloadLink string      `json:"downloadLink"`
	DriverDate   string      `json:"date"`
	DriverName   string      `json:"driverName"`
	DriverSize   string      `json:"driverSize"`
	Item         string      `json:"item"`
	Lang         string      `json:"lang"`
	Os           string      `json:"os"`
	Sha256       string      `json:"sha256"`
	SsID         interface{} `json:"ssID"`
	SvID         interface{} `json:"svID"`
	Type         string      `json:"type"`
	Vendorid     interface{} `json:"vendorID"`
	Version      string      `json:"version"`
	Updateime    string      `json:"updateTime"`
}
