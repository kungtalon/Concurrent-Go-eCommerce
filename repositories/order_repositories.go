package repositories

import (
	"gorm.io/gorm"
	"products/common"
	"products/datamodels"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (uint, error)
	Delete(uint) bool
	Update(*datamodels.Order) error
	SelectByKey(uint) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderManagerRepository struct {
	mysqlConn *gorm.DB
}

func NewOrderManagerRepository(sql *gorm.DB) IOrderRepository {
	return &OrderManagerRepository{mysqlConn: sql}
}

func (o *OrderManagerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConnGorm()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
		o.mysqlConn.AutoMigrate(&datamodels.Order{})
	}
	return nil
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (orderId uint, err error) {
	if err = o.Conn(); err != nil {
		return
	}

	result := o.mysqlConn.Create(&order)

	return order.ID, result.Error
}

func (o *OrderManagerRepository) Delete(orderId uint) bool {
	if err := o.Conn(); err != nil {
		return false
	}

	result := o.mysqlConn.Delete(&datamodels.Product{}, orderId)
	if result.Error != nil {
		return false
	}
	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) (err error) {
	if err = o.Conn(); err != nil {
		return
	}

	result := o.mysqlConn.Model(&datamodels.Product{}).Where("ID=?", order.ID).Updates(order)
	return result.Error
}

func (o *OrderManagerRepository) SelectByKey(orderId uint) (order *datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return &datamodels.Order{}, err
	}

	order = &datamodels.Order{}
	result := o.mysqlConn.First(&order, orderId)
	err = result.Error

	return
}

func (o *OrderManagerRepository) SelectAll() (orderArray []*datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return nil, err
	}

	results := o.mysqlConn.Find(&orderArray)
	err = results.Error
	return
}

func (o *OrderManagerRepository) SelectAllWithInfo() (infoMap map[int]map[string]string, err error) {
	if err = o.Conn(); err != nil {
		return nil, err
	}

	//sql := `SELECT o.ID, p.ProductName, o.OrderStatus
	//		FROM ` + o.table + ` o
	//		LEFT JOIN ` + common.PRODUCT_TABLE_NAME + ` p
	//		ON o.ProductId = p.ID`

	rows, errRows := o.mysqlConn.Table("orders").Select("orders.ID, products.product_name, orders.order_status").Joins("left join products on orders.product_id = products.ID").Rows()
	if errRows != nil {
		return nil, errRows
	}

	infoMap = common.GetResultRows(rows)
	return
}
