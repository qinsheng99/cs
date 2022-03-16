package utils

var ScoreMap map[string]string

func init() {
	ScoreMap = make(map[string]string)

	ScoreMap["AV:N"] = "Network"
	ScoreMap["AV:A"] = "Adjacent"
	ScoreMap["AV:L"] = "Local"
	ScoreMap["AV:P"] = "Physical"
	ScoreMap["AC:L"] = "Low"
	ScoreMap["AC:H"] = "High"
	ScoreMap["PR:N"] = "None"
	ScoreMap["PR:L"] = "Low"
	ScoreMap["PR:H"] = "High"
	ScoreMap["UI:N"] = "None"
	ScoreMap["UI:R"] = "Required"
	ScoreMap["S:U"] = "Unchanged"
	ScoreMap["S:C"] = "Changed"
	ScoreMap["C:N"] = "None"
	ScoreMap["C:L"] = "Low"
	ScoreMap["C:H"] = "High"
	ScoreMap["I:N"] = "None"
	ScoreMap["I:L"] = "Low"
	ScoreMap["I:H"] = "High"
	ScoreMap["A:N"] = "None"
	ScoreMap["A:L"] = "Low"
	ScoreMap["A:H"] = "High"

}
