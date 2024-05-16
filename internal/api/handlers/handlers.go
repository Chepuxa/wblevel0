package handlers

import (
	"sync"
	"wbtech/level0/internal/api/cache"
	"wbtech/level0/internal/db/repositories"

	"github.com/go-playground/validator/v10"
)

type Handlers struct {
	OrderHandler *OrderHandler
	IndexHandler *IndexHandler
}

func NewHandlers(or *repositories.OrderRepository, c *cache.Cache, wg *sync.WaitGroup, va *validator.Validate) *Handlers {
	oh := NewOrderHanlder(or, c, wg, va)
	return &Handlers{
		OrderHandler: oh,
		IndexHandler: NewIndexHandler(oh),
	}
}
