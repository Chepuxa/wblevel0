package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"wbtech/level0/internal/api/cache"
	"wbtech/level0/internal/api/models"
	"wbtech/level0/internal/customerrors"
	"wbtech/level0/internal/db/repositories"

	"github.com/go-playground/validator/v10"
)

type OrderHandler struct {
	OrderRepository *repositories.OrderRepository
	Validate        *validator.Validate
	Cache           *cache.Cache
	Wg              *sync.WaitGroup
}

func NewOrderHanlder(or *repositories.OrderRepository, ca *cache.Cache, wg *sync.WaitGroup, va *validator.Validate) *OrderHandler {
	return &OrderHandler{
		OrderRepository: or,
		Validate:        va,
		Cache:           ca,
		Wg:              wg,
	}
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	order_uid := r.URL.Query().Get("order_uuid")

	if order_uid == "" {
		customerrors.NotFoundResponse(w, r)
		return
	}

	var order models.Data

	order, ok := h.Cache.Get(order_uid)

	if !ok {

		bdorder, crudErr := h.OrderRepository.GetById(order_uid)

		if crudErr == sql.ErrNoRows {
			customerrors.NotFoundResponse(w, r)
			return
		}

		if crudErr != nil {
			customerrors.ServerErrorResponse(w, r, crudErr)
			return
		}

		order = bdorder
	}

	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	enc.Encode(order)
}
