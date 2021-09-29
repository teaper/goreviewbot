package msgc

import "strings"

// OT 话题判断
func OtMessage(msg string) (bool) {
	ots := [...]string{"屄",
		"pornhub"}
	for _, ot := range ots {
		if strings.Contains(msg, ot) {
			return true
		}
	}
	return false
}

// 复读机感叹号判断
func RepMessage(msg string) (bool)  {
	reps := [...]string{"！","!"}
	for _, ot := range reps {
		if strings.Contains(msg, ot) {
			return true
		}
	}
	return false
}
