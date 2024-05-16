package handlers

import (
	"net/http"
	"text/template"
)

type IndexHandler struct {
	orderHandler *OrderHandler
}

var tpl = template.Must(template.ParseFiles("index.html"))

func NewIndexHandler(oh *OrderHandler) *IndexHandler {
	return &IndexHandler{
		orderHandler: oh,
	}
}

func (h *IndexHandler) Handle(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}
