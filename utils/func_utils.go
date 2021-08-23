package utils

import (
	"github.com/iancoleman/strcase"
	"strings"
	"time"
)

func Ternary(value string, eq string, a string, b string) string {
	if value == eq {
		return a
	} else {
		return b
	}
}

func TernaryBool(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func FetchListFromString(subjectString string, demiliter string) []string {
	if subjectString == "" {
		return []string{}
	}
	return strings.Split(subjectString, demiliter)
}
func FetchInQueryStringFromArray(array []string) string {
	resultString := ""
	for _, val := range array {
		resultString = resultString + "'" + val + "'" + ","
	}
	return strings.Trim(resultString, ",")
}

func FetchInQueryStringFromMapKeys(inputMap map[string]int) string {
	resultString := ""
	for key := range inputMap {
		resultString = resultString + "'" + key + "'" + ","
	}
	return strings.Trim(resultString, ",")
}

func ExistsInArray(objArray []string, target string) bool {
	for _, val := range objArray {
		if val == target {
			return true
		}
	}
	return false
}

func ValidateSize(size string) bool {
	sizeObj := StringToSize(size)
	for _, validSize := range sizeList {
		if sizeObj == validSize {
			return true
		}
	}
	return false
}

func ValidateStatus(statusString string) bool{
	for _, status := range statusList{
		if status == statusString{
			return true
		}
	}
	return false
}

func SizeToString(size SIZE) string {
	return string(size)
}

func StringToSize(size string) SIZE {
	return SIZE(size)
}

func StringFromArray(array []string, character string) string {
	result := ""
	for _, val := range array {
		result = result + (val + character)
	}
	result = strings.Trim(result, character)
	return result
}

func Unique(e []string) []string {
	r := []string{}
	for _, s := range e {
		if !Contains(r[:], s) {
			r = append(r, s)
		}
	}
	return r
}

func Contains(e []string, c string) bool {
	for _, s := range e {
		if s == c {
			return true
		}
	}
	return false
}

func ToScreamingSnakeArray(str []string) []string {
	res := []string{}
	for _, s := range str {
		res = append(res, strcase.ToScreamingSnake(s))
	}
	return res
}

func ParseStringToTime(timeString string) (time.Time, error) {
	timeValue, err := time.Parse("2006-01-02 15:04:05", timeString)
	if err!=nil{
		return time.Time{}, err
	}
	//durationToAdd, _ := time.ParseDuration("5h30m")
	//timeValue.Add(durationToAdd)
	return timeValue, nil
}

func RemoveAll(s []string, t []string) []string {
	r := []string{}
	for _, e := range s {
		if !Contains(t, e) {
			r = append(r, e)
		}
	}
	return r
}
