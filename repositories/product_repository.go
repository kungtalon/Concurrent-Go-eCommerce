package repositories

import (
	"database/sql"
	"products/common"
	"products/datamodels"
	"strconv"
)

// develop the interface
// implement the interface

type IProduct interface {
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}

func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

func (p *ProductManager) Insert(product *datamodels.Product) (productID int64, err error) {
	// check connection
	if err = p.Conn(); err != nil {
		return
	}

	// prepare formatted sql statement
	sql := `
		INSERT ` + p.table + ` 
		SET productName=?,productNum=?,productImage=?,productUrl=?`
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nil {
		return 0, errSql
	}

	// pass in arguments in sql
	result, errStmt := stmt.Exec(
		product.ProductName,
		product.ProductNum,
		product.ProductImage,
		product.ProductUrl,
	)
	if errStmt != nil {
		return 0, errStmt
	}
	return result.LastInsertId()
}

func (p *ProductManager) Delete(productID int64) bool {
	// check connection
	if err := p.Conn(); err != nil {
		return false
	}

	sql := "delete from " + p.table + " where ID=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}

	_, err = stmt.Exec(productID)
	if err != nil {
		return false
	}
	return true
}

func (p *ProductManager) Update(product *datamodels.Product) (err error) {
	if err = p.Conn(); err != nil {
		return err
	}

	sql := `
		UPDATE ` + p.table + ` 
		SET productNum=?,productName=?,productImage=?,productUrl=? 
		WHERE ID=` + strconv.FormatInt(product.ID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		product.ProductNum,
		product.ProductName,
		product.ProductImage,
		product.ProductUrl,
	)
	if err != nil {
		return err
	}

	return
}

// select a product by its ID
func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
	// check connection
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}

	// prepare formatted sql statement
	sql := `
		SELECT * 
		FROM  ` + p.table + `
		WHERE ID=` + strconv.FormatInt(productID, 10)

	row, errRow := p.mysqlConn.Query(sql)
	// be careful when we are going to copy data from sql rows
	if errRow != nil {
		return &datamodels.Product{}, errRow
	}
	defer row.Close()

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}

	common.DataToStructByTagSql(result, productResult)
	return
}

func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, err error) {
	if err = p.Conn(); err != nil {
		return nil, err
	}

	sql := `SELECT * FROM ` + p.table
	rows, err := p.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := common.GetResultRows(rows)
	if len(results) == 0 {
		return nil, nil
	}

	for _, v := range results {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}

	return
}
