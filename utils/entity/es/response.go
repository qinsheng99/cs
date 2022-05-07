package es

type DriverEsData struct {
	Id           int64  `json:"id"`
	Architecture string `json:"architecture"`
	BoardModel   string `json:"boardModel"`
	ChipModel    string `json:"chipModel"`
	ChipVendor   string `json:"chipVendor"`
	Deviceid     string `json:"deviceID"`
	DownloadLink string `json:"downloadLink"`
	DriverDate   string `json:"driverDate"`
	DriverName   string `json:"driverName"`
	DriverSize   string `json:"driverSize"`
	Item         string `json:"item"`
	Lang         string `json:"lang"`
	Os           string `json:"os"`
	Sha256       string `json:"sha256"`
	SsID         string `json:"ssID"`
	SvID         string `json:"svID"`
	Type         string `json:"type"`
	Vendorid     string `json:"vendorID"`
	Version      string `json:"version"`
	Updateime    string `json:"updateTime"`
}

type DriverEsDataResp struct {
	Total  int64          `json:"total"`
	Driver []DriverEsData `json:"driver"`
}

type Gotty struct {
	EventType     int64  `json:"eventType"`
	Input         string `json:"input"`
	Instance      string `json:"instance"`
	OperationTime string `json:"operationTime"`
	Output        string `json:"output"`
	Ps            string `json:"ps"`
}

type GottyEsDataResp struct {
	Total  int64   `json:"total"`
	Driver []Gotty `json:"driver"`
}
