package repositories

import "database/sql"

type Repositories struct {
	OrderRepository *OrderRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		OrderRepository: NewOrderRepository(db),
	}
}
