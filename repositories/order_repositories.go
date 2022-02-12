package repositories

import (
	"database/sql"
	"products/common"
	"products/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func NewOrderManagerRepository(table string, sql *sql.DB) IOrderRepository {
	return &OrderManagerRepository{table: table, mysqlConn: sql}
}

func (o *OrderManagerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = common.ORDER_TABLE_NAME
	}
	return nil
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (orderId int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}

	sql := `
		INSERT` + o.table + `
		SET UserID=?, ProductID=?, OrderStatus=?`
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return orderId, errStmt
	}
	result, errResult := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if errResult != nil {
		return orderId, errResult
	}

	return result.LastInsertId()
}

func (o *OrderManagerRepository) Delete(orderId int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}

	sql := `DELETE FROM ` + o.table + `
			WHERE ID=?`
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return false
	}

	_, err := stmt.Exec(orderId)
	if err != nil {
		return false
	}
	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) (err error) {
	if err = o.Conn(); err != nil {
		return
	}

	sql := `UPDATE ` + o.table + `
			SET UserId=?, ProductId=?, OrderStatus=?
			WHERE ID=` + strconv.FormatInt(order.ID, 10)
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return errStmt
	}

	_, err = stmt.Exec(stmt)
	return
}

func (o *OrderManagerRepository) SelectByKey(orderId int64) (order *datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return &datamodels.Order{}, err
	}

	sql := `SELECT * FROM ` + o.table + `WHERE ID=` + strconv.FormatInt(orderId, 10)
	row, errRow := o.mysqlConn.Query(sql)
	if errRow != nil {
		return &datamodels.Order{}, errRow
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, nil
	}

	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

func (o *OrderManagerRepository) SelectAll() (orderArray []*datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return nil, err
	}

	sql := `SELECT * FROM ` + o.table
	rows, errRow := o.mysqlConn.Query(sql)
	if errRow != nil {
		return nil, errRow
	}

	results := common.GetResultRows(rows)
	if len(results) == 0 {
		return nil, nil
	}

	for _, result := range results {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(result, order)
		orderArray = append(orderArray, order)
	}
	return
}

func (o *OrderManagerRepository) SelectAllWithInfo() (infoMap map[int]map[string]string, err error) {
	if err = o.Conn(); err != nil {
		return nil, err
	}

	sql := `SELECT o.ID, p.ProductName, o.OrderStatus
			FROM ` + o.table + ` o
			LEFT JOIN ` + common.PRODUCT_TABLE_NAME + ` p
			ON o.ProductId = p.ID`

	rows, errRows := o.mysqlConn.Query(sql)
	if errRows != nil {
		return nil, errRows
	}

	infoMap = common.GetResultRows(rows)
	return
}
