package common

import (
	"bufio"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	gomysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"strings"
)

func ReadDBConn(path string) (authen []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		contents := strings.Split(scanner.Text(), " ")
		authen = append(authen, contents...)
	}
	return
}

// create a mysql connection
func NewMysqlConn() (db *sql.DB, err error) {
	authentication, err := ReadDBConn("common/DBCON")
	if err != nil {
		return
	}
	db, err = sql.Open("mysql", authentication[0]+":"+authentication[1]+"@tcp(localhost:3306)/lightning?charset=utf8")
	return
}

func NewMysqlConnGorm() (db *gorm.DB, err error) {
	//mysqlDB, err := NewMysqlConn()
	authentication, err := ReadDBConn("common/DBCON")
	if err != nil {
		return
	}
	dsn := authentication[0] + ":" + authentication[1] + "@tcp(localhost:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(gomysql.Open(dsn), &gorm.Config{})
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
	// return all columns in the format of a map
	columns, _ := rows.Columns()
	// use []byte to represent a row
	vals := make([][]byte, len(columns))
	// a row of filling data
	scans := make([]interface{}, len(columns))
	// temporarily we scan row values into []byte
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		// read row data
		rows.Scan(scans...)
		// create a map for this row
		row := make(map[string]string)
		// copy data in vals to row
		for k, v := range vals {
			key := columns[k]
			// convert []byte to string
			row[key] = string(v)
		}
		// put into the map that will be returned
		result[i] = row
		i++
	}
	return result
}
