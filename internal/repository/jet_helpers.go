package repository

import "github.com/go-jet/jet/v2/mysql"

func stringExprOrNull(value *string) mysql.StringExpression {
	if value == nil {
		return mysql.StringExp(mysql.NULL)
	}
	return mysql.String(*value)
}
