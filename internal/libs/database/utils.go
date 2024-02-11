package database

import "strings"

type SqlQuery struct {
	Count      string   `json:"count"`
	GroupBy    []string `json:"group_by"`
	OrderBy    string   `json:"order_by"`
	TableName  string   `json:"table_name"`
	Conditions []Condition
}

type Condition struct {
	Field string
	Value string
	Op    string
}

func ParseQueryString(queryString string) (*SqlQuery, error) {
	queryParts := strings.Split(queryString, " ")
	tableName := queryParts[1]

	var sqlQuery SqlQuery
	sqlQuery.TableName = tableName

	for i, part := range queryParts {
		switch part {
		case "SELECT":
			if i+1 < len(queryParts) && queryParts[i+1] == "COUNT" {
				sqlQuery.Count = queryParts[i+2]
				//i += 2
			}
		case "FROM":
			continue
		case "GROUP":
			for j := i + 1; j < len(queryParts); j++ {
				if queryParts[j] == "BY" {
					sqlQuery.GroupBy = queryParts[i+2 : j]
					break
				}
			}
		case "ORDER":
			for j := i + 1; j < len(queryParts); j++ {
				if queryParts[j] == "BY" {
					sqlQuery.OrderBy = queryParts[i+2]
					break
				}
			}
		case "WHERE":
			for j := i + 1; j < len(queryParts); j++ {
				if queryParts[j] == "AND" || queryParts[j] == "OR" {
					break
				}
				condition := Condition{
					Field: queryParts[j],
					Op:    "=",
				}
				if j+1 < len(queryParts) && queryParts[j+1] == ">" {
					condition.Op = ">"
					condition.Value = queryParts[j+2]
					j += 2
				} else {
					condition.Value = queryParts[j+1]
					j++
				}
				sqlQuery.Conditions = append(sqlQuery.Conditions, condition)
			}
		}
	}

	return &sqlQuery, nil
}
