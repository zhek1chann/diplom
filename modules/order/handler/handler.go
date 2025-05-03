package handler

type OrderHandler struct {
	service IOrderService
}

func NewHandler(service IOrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

type IOrderService interface {
	ICreateOrderService
}
