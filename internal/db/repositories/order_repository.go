package repositories

import (
	"database/sql"
	"wbtech/level0/internal/api/models"
)

type OrderRepositoryInterface interface {
	Create(*models.Order) error
	GetById(int64) (models.Order, error)
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(data *models.Data) error {

	_, err := r.db.Exec("INSERT INTO orders VALUES($1, $2)", data.OrderUID, data)

	return err
}

func (r *OrderRepository) GetById(order_uuid string) (models.Data, error) {
	data := new(models.Data)

	err := r.db.QueryRow("SELECT data FROM orders WHERE order_uid = $1;", order_uuid).Scan(&data)

	return *data, err
}

func (r *OrderRepository) GetAll() ([]models.Data, error) {
	orders := []models.Data{}

	rows, queryErr := r.db.Query("SELECT data FROM orders")

	if queryErr != nil {
		return nil, queryErr
	}

	defer rows.Close()

	for rows.Next() {
		var data models.Data

		scanErr := rows.Scan(&data)

		if scanErr != nil {
			return nil, scanErr
		}

		orders = append(orders, data)

	}

	return orders, nil
}
