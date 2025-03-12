package model

type PageCountInput struct {
	PageSize int `json:"pageSize"`
}

type PageCountResponse struct {
	PageCount int `json:"pageCount"`
}
