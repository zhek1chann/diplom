package model

type ErrorResponse struct{
	Err error `json:"error"`
}