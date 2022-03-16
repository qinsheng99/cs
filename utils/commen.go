package utils

import (
	"errors"
	"strings"
	"time"

	_const "cve-sa-backend/utils/const"

	"gorm.io/gorm"
)

func GetCurTime() string {
	return time.Now().Format(_const.Format)
}

func TrimString(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, "\t", "", -1)
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
	vectorValue := VectorValue{}
	str = GetTextTrim(str)
	arr := strings.Split(str, "/")
	var innerArr []string
	for _, v := range arr {
		v = strings.TrimSpace(v)
		v = strings.Replace(v, "：", ":", -1)
		v = strings.Replace(v, " ", "", -1)
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