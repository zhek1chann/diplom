package handler

type UserHandler struct {
	service IUserService
}

func NewHandler(service IUserService) *UserHandler {
	return &UserHandler{service: service}
}

type IUserService interface {
	IAddressService
}
