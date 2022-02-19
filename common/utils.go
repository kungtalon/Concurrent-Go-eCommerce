package common

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

// map sql query results to struct according to the struct tag
func DataToStructByTagSql(data map[string]string, obj interface{}) {
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValue.NumField(); i++ {
		// get the value of sql
		value := data[objValue.Type().Field(i).Tag.Get("sql")]
		// get the corresponding field name
		name := objValue.Type().Field(i).Name
		// get the target data type
		structFieldType := objValue.Field(i).Type()
		// get data type
		val := reflect.ValueOf(value)
		var err error
		if structFieldType != val.Type() {
			// datatype conversion
			val, err = TypeConversion(value, structFieldType.Name())
			if err != nil {

			}
		}
		// set the value type, here we are using pointers!
		objValue.FieldByName(name).Set(val)
	}
}

// datatype conversion
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......

	return reflect.ValueOf(value), errors.New("unknown type: " + ntype)
}