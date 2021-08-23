package utils

import "strings"

func CreateUpdateFieldStatement(fieldList, valueList []string) string {
	statement := ""
	if len(fieldList) != len(valueList) {
		return ""
	}
	for index := 0; index < len(fieldList); index++ {
		statement += fieldList[index] + "='" + valueList[index] + "',"
	}
	statement = strings.Trim(statement, ",")
	return statement
}
