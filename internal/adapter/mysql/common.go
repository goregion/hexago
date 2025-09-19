package adapter_mysql

import _ "github.com/go-sql-driver/mysql"

func makeOHLCTableName(timeframeName string) string {
	return "`ohlc_" + timeframeName + "`"
}
