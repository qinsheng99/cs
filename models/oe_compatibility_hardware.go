package models

type OeCompatibilityHardware struct {
	Id                  int64  `gorm:"column:id" json:"id"`
	Architecture        string `gorm:"column:architecture" json:"architecture"`
	BiosUefi            string `gorm:"column:bios_uefi" json:"biosUefi"`
	CertificationAddr   string `gorm:"column:certification_addr" json:"certificationAddr"`
	CertificationTime   string `gorm:"column:certification_time" json:"certificationTime"`
	CommitID            string `gorm:"column:commitid" json:"commitID"`
	ComputerType        string `gorm:"column:computer_type" json:"computerType"`
	Cpu                 string `gorm:"column:cpu" json:"cpu"`
	Date                string `gorm:"column:date" json:"date"`
	FriendlyLink        string `gorm:"column:friendly_link" json:"friendlyLink"`
	HardDiskDrive       string `gorm:"column:hard_disk_drive" json:"hardDiskDrive"`
	HardwareFactory     string `gorm:"column:hardware_factory" json:"hardwareFactory"`
	HardwareModel       string `gorm:"column:hardware_model" json:"hardwareModel"`
	HostBusAdapter      string `gorm:"column:host_bus_adapter" json:"hostBusAdapter"`
	Lang                string `gorm:"column:lang" json:"lang"`
	MotherBoardRevision string `gorm:"column:mother_board_revision" json:"motherBoardRevision"`
	OsVersion           string `gorm:"column:os_version" json:"osVersion"`
	PortsBusTypes       string `gorm:"column:ports_bus_types" json:"portsBusTypes"`
	ProductInformation  string `gorm:"column:product_information" json:"productInformation"`
	Ram                 string `gorm:"column:ram" json:"ram"`
	VideoAdapter        string `gorm:"column:video_adapter" json:"videoAdapter"`
	Updateime           string `gorm:"column:update_time" json:"updateTime"`
}

func (o *OeCompatibilityHardware) TableName() string {
	return "oe_compatibility_hardware"
}

type ROeCompatibilityHardware struct {
	OeCompatibilityHardware
	Updateime string `json:"updateTime"`
}
