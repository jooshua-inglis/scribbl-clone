package data

import "fmt"

func dynamicUpdateSet(values []string, updateSet map[string]any) string {
	query := ""
	var insertComma = false
	for i := range values {
		value := values[i]
		if _, ok := updateSet[value]; ok {
			if insertComma {
				query += ","
			}
			query += fmt.Sprintf(" %s = :%s ", toSnake(value), value)
			insertComma = true
		}
	}
	return query
}
