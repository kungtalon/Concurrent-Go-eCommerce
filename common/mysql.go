package common

import (
	"database/sql"
)

// create a mysql connection
func NewMysqlConn() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "root:imooc@tcp(127.0.0.1:3306)/imooc?charset=utf8")
	return
}

// return one row from query results
func GetResultRow(rows *sql.Rows) map[string]string {
	columns, _ := rows.Columns()
	// to make this reading function more generic
	// assume that the target data type is unknown
	// so we need to use bytes to pull out the sqlRow into Go
	// and then convert bytes to string
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[string]string)
	for rows.Next() {
		//将行数据保存到record字典
		rows.Scan(scanArgs...)
		for i, v := range values {
			if v != nil {
				//fmt.Println(reflect.TypeOf(col))
				record[columns[i]] = string(v)
			}
		}
	}
	return record
}

// get all
func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	// return all columns
	columns, _ := rows.Columns()
	// use []byte to represent a row
	vals := make([][]byte, len(columns))
	// a row of filling data
	scans := make([]interface{}, len(columns))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := columns[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	return result
}
