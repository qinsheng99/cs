package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"

	_const "cve-sa-backend/utils/const"
	cveSa "cve-sa-backend/utils/entity/cve_sa"

	"gorm.io/gorm"
)

func GetCurTime() string {
	return time.Now().Format(_const.Format)
}

//TrimString Remove the \n \r \t " " in the string
func TrimString(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, "\t", "", -1)
	return str
}

//TrimStringNR Remove the \n \r in the string
func TrimStringNR(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}

func RemoveChSymbols(s string) string {
	s = strings.Replace(s, "：", ":", -1)
	s = strings.Replace(s, " ", "", -1)
	return s
}

func ErrorNotFound(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	return false
}

func GetTextTrim(str string) string {
	str = strings.Replace(str, "\t", " ", -1)
	str = strings.Replace(str, "\n", " ", -1)
	str = strings.Replace(str, "\r", " ", -1)
	str = strings.Replace(str, "\f", " ", -1)
	str = strings.TrimSpace(str)
	return str
}

func GetVectorArr(str string) VectorValue {
	if str[:1] == "(" {
		str, _ = InterceptString(str, "(", ")")
	}

	vectorValue := VectorValue{}
	str = GetTextTrim(str)
	arr := strings.Split(str, "/")
	var innerArr []string
	for _, v := range arr {
		v = strings.TrimSpace(v)
		v = RemoveChSymbols(v)
		innerArr = strings.Split(v, ":")
		if innerArr[0] == "AV" {
			val, ok := ScoreMap[v]
			if ok {
				vectorValue.AV = val
			}
		}
		if innerArr[0] == "AC" {
			val, ok := ScoreMap[v]
			if ok {
				vectorValue.AC = val
			}
		}
		if innerArr[0] == "PR" {
			val, ok := ScoreMap[v]
			if ok {
				vectorValue.PR = val
			}
		}
		if innerArr[0] == "UI" {
			val, ok := ScoreMap[v]
			if ok {
				vectorValue.UI = val
			}
		}
		if innerArr[0] == "S" {
			val, ok := ScoreMap[v]
			if ok {
				vectorValue.S = val
			}
		}
		if innerArr[0] == "C" {
			val, ok := ScoreMap[v]
			if ok {
				vectorValue.C = val
			}
		}
		if innerArr[0] == "I" {
			val, ok := ScoreMap[v]
			if ok {
				vectorValue.I = val
			}
		}
		if innerArr[0] == "A" {
			val, ok := ScoreMap[v]
			if ok {
				vectorValue.A = val
			}
		}

	}
	return vectorValue
}

type VectorValue struct {
	AV string
	AC string
	PR string
	UI string
	S  string
	C  string
	I  string
	A  string
}

func InterfaceToString(i interface{}) string {
	switch x := i.(type) {
	case float64:
		return strconv.Itoa(int(x))
	case float32:
		return strconv.Itoa(int(x))
	case int:
		return strconv.Itoa(x)
	case int8:
		return strconv.Itoa(int(x))
	case int16:
		return strconv.Itoa(int(x))
	case int32:
		return strconv.Itoa(int(x))
	case int64:
		return strconv.Itoa(int(x))
	case string:
		return x
	default:
		return ""
	}
}

func GetPage(req cveSa.Pages) (int, int) {
	var page, size int
	if req.Page == 0 {
		page = _const.Page
	} else {
		page = req.Page
	}
	if req.Size == 0 {
		size = _const.Size
	} else {
		size = req.Size
	}
	return page, size
}
