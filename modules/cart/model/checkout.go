package model

const (
	PaymentStatusSuccess = "success"
	PaymentStatusFailed  = "failed"
)

type CheckoutResponse struct {
	CheckoutURL    string
	PaymentID      string
	ResponseStatus string
}

type PaymentOrder struct {
	ID   string
	Cart Cart
	CheckoutResponse
}

type CommitCheckout struct {
	OrderID       string
	PaymentStatus string
}
