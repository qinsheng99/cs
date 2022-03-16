package cveSa

type RequestData struct {
	KeyWord     string `json:"keyword"`
	Type        string `json:"type"`
	Year        string `json:"year"`
	Status      string `json:"status"`
	PackageName string `json:"packageName"`
	Pages       Pages  `json:"pages"`
}

type Pages struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type OeCompSearchRequest struct {
	Os           string `json:"os"`
	Architecture string `json:"architecture"`
	KeyWord      string `json:"keyword"`
	Lang         string `json:"lang"`
	Pages        Pages  `json:"pages"`
}

type SyncRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DeleteSAID  string `json:"deleteSAID"`
	DeleteCVEID string `json:"deleteCVEID"`
	PackageName string `json:"packageName"`
	CveNo       string `json:"cveNo"`
	SaNo        string `json:"saNo"`
}
